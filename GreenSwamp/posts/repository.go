package posts

import (
	"context"
	"errors"
	"fmt"

	"gorm.io/gorm"
)

type Repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{db: db}
}

func AutoMigrate(db *gorm.DB) error {
	if err := db.AutoMigrate(&User{}, &Post{}, &Interaction{}, &Event{}, &Tag{}, &PostTag{}); err != nil {
		return err
	}
	return ensureTrendingView(db)
}

func (r *Repository) ListFeed(ctx context.Context) ([]Post, error) {
	var posts []Post
	err := r.db.WithContext(ctx).
		Preload("User").
		Preload("Event").
		Preload("Tags").
		Order("posts.created_at DESC").
		Find(&posts).Error
	if err != nil {
		return nil, err
	}
	return r.attachInteractionCounts(ctx, posts)
}

func (r *Repository) ListPostsByTag(ctx context.Context, tag string) ([]Post, error) {
	tag = normalizeTagName(tag)

	var posts []Post
	err := r.db.WithContext(ctx).
		Model(&Post{}).
		Select("posts.*").
		Joins("JOIN post_tags pt ON pt.post_id = posts.post_id").
		Joins("JOIN tags t ON t.tag_id = pt.tag_id").
		Where("LOWER(t.tag_name) = LOWER(?)", tag).
		Group("posts.post_id").
		Preload("User").
		Preload("Event").
		Preload("Tags").
		Order("posts.created_at DESC").
		Find(&posts).Error
	if err != nil {
		return nil, err
	}

	return r.attachInteractionCounts(ctx, posts)
}

func (r *Repository) GetPostByID(ctx context.Context, postID uint) (*Post, error) {
	if err := r.db.WithContext(ctx).
		Model(&Post{}).
		Where("post_id = ?", postID).
		Update("view_count", gorm.Expr("view_count + 1")).Error; err != nil {
		return nil, err
	}

	var post Post
	err := r.db.WithContext(ctx).
		Preload("User").
		Preload("Event").
		Preload("Tags").
		First(&post, "post_id = ?", postID).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}
	if err != nil {
		return nil, err
	}

	posts, err := r.attachInteractionCounts(ctx, []Post{post})
	if err != nil {
		return nil, err
	}
	return &posts[0], nil
}

func (r *Repository) ListProfileByUsername(ctx context.Context, username string) (*User, []Post, error) {
	var user User
	if err := r.db.WithContext(ctx).First(&user, "username = ?", username).Error; err != nil {
		return nil, nil, err
	}

	var posts []Post
	if err := r.db.WithContext(ctx).
		Preload("User").
		Preload("Event").
		Preload("Tags").
		Where("user_id = ?", user.UserID).
		Order("posts.created_at DESC").
		Find(&posts).Error; err != nil {
		return nil, nil, err
	}

	posts, err := r.attachInteractionCounts(ctx, posts)
	if err != nil {
		return nil, nil, err
	}

	return &user, posts, nil
}

func (r *Repository) ListTrendingPonds(ctx context.Context, limit int) ([]TrendingPond, error) {
	if limit <= 0 {
		limit = 10
	}

	var tags []TrendingPond
	query := r.db.WithContext(ctx).
		Raw(`SELECT tag_id, tag_name, recent_posts FROM trending_ponds ORDER BY recent_posts DESC, tag_name ASC LIMIT ?`, limit)

	if err := query.Scan(&tags).Error; err != nil {
		return nil, fmt.Errorf("list trending ponds: %w", err)
	}
	return tags, nil
}

func (r *Repository) attachInteractionCounts(ctx context.Context, posts []Post) ([]Post, error) {
	if len(posts) == 0 {
		return posts, nil
	}

	type countRow struct {
		PostID          uint
		InteractionType string
		Count           int64
	}

	var rows []countRow
	if err := r.db.WithContext(ctx).
		Model(&Interaction{}).
		Select("post_id, interaction_type, COUNT(*) as count").
		Where("post_id IN ?", extractPostIDs(posts)).
		Group("post_id, interaction_type").
		Scan(&rows).Error; err != nil {
		return nil, err
	}

	for i := range posts {
		var total, likes, reribbs int64
		for _, row := range rows {
			if row.PostID != posts[i].PostID {
				continue
			}
			total += row.Count
			switch row.InteractionType {
			case "like":
				likes = row.Count
			case "reribb":
				reribbs = row.Count
			}
		}
		posts[i].InteractionCount = total
		posts[i].LikeCount = likes
		posts[i].ReribbCount = reribbs
	}

	return posts, nil
}

func extractPostIDs(posts []Post) []uint {
	ids := make([]uint, 0, len(posts))
	for _, p := range posts {
		ids = append(ids, p.PostID)
	}
	return ids
}

func (r *Repository) CreatePost(ctx context.Context, userID uint, content string, postType string) (*Post, error) {
	if content == "" {
		return nil, errors.New("content cannot be empty")
	}

	tx := r.db.WithContext(ctx).Begin()
	post := &Post{
		UserID:   userID,
		Content:  content,
		PostType: postType,
	}

	if err := tx.Create(post).Error; err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("failed to create post: %w", err)
	}

	matches := hashtagRe.FindAllStringSubmatch(content, -1)
	seen := make(map[string]struct{})
	for _, m := range matches {
		if len(m) < 3 {
			continue
		}
		tagName := normalizeTagName(m[2])
		if tagName == "" {
			continue
		}
		if _, ok := seen[tagName]; ok {
			continue
		}
		seen[tagName] = struct{}{}

		var t Tag
		err := tx.Where("LOWER(tag_name) = ?", tagName).First(&t).Error
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				t = Tag{TagName: tagName}
				if err := tx.Create(&t).Error; err != nil {
					tx.Rollback()
					return nil, fmt.Errorf("failed to create tag: %w", err)
				}
			} else {
				tx.Rollback()
				return nil, fmt.Errorf("find tag: %w", err)
			}
		}

		pt := PostTag{PostID: post.PostID, TagID: t.TagID}
		if err := tx.Create(&pt).Error; err != nil {
			tx.Rollback()
			return nil, fmt.Errorf("failed to attach tag to post: %w", err)
		}
	}

	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("commit post transaction: %w", err)
	}

	if err := r.db.WithContext(ctx).Preload("User").Preload("Tags").First(post, post.PostID).Error; err != nil {
		return nil, fmt.Errorf("failed to load created post: %w", err)
	}

	return post, nil
}

func (r *Repository) UpdatePostMedia(ctx context.Context, postID uint, mediaURL, mediaType, altText *string) error {
	return r.db.WithContext(ctx).
		Model(&Post{}).
		Where("post_id = ?", postID).
		Updates(map[string]interface{}{
			"media_url":  mediaURL,
			"media_type": mediaType,
			"alt_text":   altText,
		}).Error
}

func (r *Repository) ToggleInteraction(ctx context.Context, userID, postID uint, interactionType string) (created bool, count int64, err error) {
	var existing Interaction
	err = r.db.WithContext(ctx).
		Where("user_id = ? AND post_id = ? AND interaction_type = ?", userID, postID, interactionType).
		First(&existing).Error

	if err == nil {
		if delErr := r.db.WithContext(ctx).Delete(&existing).Error; delErr != nil {
			return false, 0, fmt.Errorf("delete interaction: %w", delErr)
		}
		count, err = r.CountByType(ctx, postID, interactionType)
		return false, count, err
	}

	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return false, 0, fmt.Errorf("find interaction: %w", err)
	}

	interaction := Interaction{
		UserID:          userID,
		PostID:          postID,
		InteractionType: interactionType,
	}
	if createErr := r.db.WithContext(ctx).Create(&interaction).Error; createErr != nil {
		return false, 0, fmt.Errorf("create interaction: %w", createErr)
	}

	count, err = r.CountByType(ctx, postID, interactionType)
	return true, count, err
}

func (r *Repository) GetUserInteractions(ctx context.Context, userID uint, postIDs []uint) (map[uint][]string, error) {
	if len(postIDs) == 0 {
		return nil, nil
	}

	var rows []Interaction
	if err := r.db.WithContext(ctx).
		Where("user_id = ? AND post_id IN ?", userID, postIDs).
		Find(&rows).Error; err != nil {
		return nil, fmt.Errorf("get user interactions: %w", err)
	}

	result := make(map[uint][]string, len(postIDs))
	for _, row := range rows {
		result[row.PostID] = append(result[row.PostID], row.InteractionType)
	}
	return result, nil
}

func (r *Repository) CountByType(ctx context.Context, postID uint, interactionType string) (int64, error) {
	var count int64
	if err := r.db.WithContext(ctx).
		Model(&Interaction{}).
		Where("post_id = ? AND interaction_type = ?", postID, interactionType).
		Count(&count).Error; err != nil {
		return 0, fmt.Errorf("count interactions: %w", err)
	}
	return count, nil
}

func (r *Repository) CreateComment(ctx context.Context, userID, postID uint, content string) (*Interaction, error) {
	comment := &Interaction{
		UserID:          userID,
		PostID:          postID,
		InteractionType: "comment",
		Content:         &content,
	}
	if err := r.db.WithContext(ctx).Create(comment).Error; err != nil {
		return nil, fmt.Errorf("create comment: %w", err)
	}
	return comment, nil
}

func (r *Repository) ListCommentsByPost(ctx context.Context, postID uint) ([]Interaction, error) {
	var comments []Interaction
	if err := r.db.WithContext(ctx).
		Where("post_id = ? AND interaction_type = ?", postID, "comment").
		Order("created_at DESC").
		Find(&comments).Error; err != nil {
		return nil, fmt.Errorf("list comments: %w", err)
	}
	return comments, nil
}

func (r *Repository) GetUserByID(ctx context.Context, userID uint) (*User, error) {
	var u User
	if err := r.db.WithContext(ctx).First(&u, "user_id = ?", userID).Error; err != nil {
		return nil, err
	}
	return &u, nil
}

func ensureTrendingView(db *gorm.DB) error {
	const q = `
	CREATE VIEW IF NOT EXISTS trending_ponds AS
	SELECT t.tag_id, t.tag_name, COUNT(pt.post_id) AS recent_posts
	FROM tags t
	JOIN post_tags pt ON t.tag_id = pt.tag_id
	JOIN posts p ON pt.post_id = p.post_id
	WHERE p.created_at > datetime('now', '-7 days')
	GROUP BY t.tag_id, t.tag_name
	`
	db.Exec("DROP VIEW IF EXISTS trending_ponds")
	return db.Exec(q).Error
}