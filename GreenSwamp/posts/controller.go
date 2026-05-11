package posts

import (
	"context"
	"encoding/json"
	"html/template"
	"log"
	"net/http"
	"path"
	"strconv"
	"strings"

	csrf "github.com/gorilla/csrf"
	account "github.com/venexene/serv-prog-go/greenswamp/account"
	"gorm.io/gorm"
)

type Controller struct {
	repo     *Repository
	authRepo *account.Repository
	tmpl     *template.Template
	basePath string
	logger   *log.Logger
}

func (c *Controller) handleIndex(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != c.basePath && r.URL.Path != c.basePath+"/" {
		http.NotFound(w, r)
		return
	}
	http.Redirect(w, r, c.basePath+"/feed", http.StatusFound)
}

func (c *Controller) handlePondsIndex(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != c.basePath+"/ponds" {
		http.NotFound(w, r)
		return
	}
	http.Redirect(w, r, c.basePath+"/feed", http.StatusFound)
}

func (c *Controller) handleFeed(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != c.basePath+"/feed" && r.URL.Path != c.basePath+"/feed/" {
		http.NotFound(w, r)
		return
	}

	posts, err := c.repo.ListFeed(r.Context())
	if err != nil {
		c.renderError(w, http.StatusInternalServerError, "could not load feed")
		return
	}

	data := FeedPage{
		Title:       "Greenswamp · Feed",
		BasePath:    c.basePath,
		Items:       c.buildItems(posts),
		Trending:    c.loadTrending(r.Context()),
		CurrentUser: c.currentUser(r),
		CSRFToken:   csrf.Token(r),
	}

	c.execute(w, "feed.html", data)
}

func (c *Controller) handlePond(w http.ResponseWriter, r *http.Request) {
	tag := strings.TrimPrefix(r.URL.Path, c.basePath+"/ponds/")
	tag = path.Clean("/" + tag)
	tag = strings.TrimPrefix(tag, "/")

	if tag == "" || tag == "." || strings.Contains(tag, "/") {
		http.Redirect(w, r, c.basePath+"/feed", http.StatusFound)
		return
	}

	tag = normalizeTagName(tag)
	if tag == "" {
		http.Redirect(w, r, c.basePath+"/feed", http.StatusFound)
		return
	}

	posts, err := c.repo.ListPostsByTag(r.Context(), tag)
	if err != nil {
		c.renderError(w, http.StatusInternalServerError, "could not load topic")
		return
	}

	data := FeedPage{
		Title:       "Greenswamp · #" + tag,
		BasePath:    c.basePath,
		Tag:         tag,
		Items:       c.buildItems(posts),
		Trending:    c.loadTrending(r.Context()),
		CurrentUser: c.currentUser(r),
	}

	c.execute(w, "feed.html", data)
}

func (c *Controller) handleProfile(w http.ResponseWriter, r *http.Request) {
	username := strings.TrimPrefix(r.URL.Path, c.basePath+"/profile/")
	username = path.Clean("/" + username)
	username = strings.TrimPrefix(username, "/")

	if username == "" || username == "." || strings.Contains(username, "/") {
		http.NotFound(w, r)
		return
	}

	user, posts, err := c.repo.ListProfileByUsername(r.Context(), username)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			http.NotFound(w, r)
			return
		}
		c.renderError(w, http.StatusInternalServerError, "could not load profile")
		return
	}

	profile := ProfilePage{
		User:      *user,
		Posts:     c.buildItems(posts),
		PostCount: len(posts),
		Avatar:    avatarOrFallback(user.AvatarURL),
		Bio:       bioOrEmpty(user.Bio),
	}

	data := ProfilePageData{
		Title:       user.DisplayName + " · Profile",
		BasePath:    c.basePath,
		Profile:     profile,
		Trending:    c.loadTrending(r.Context()),
		CurrentUser: c.currentUser(r),
	}

	c.execute(w, "profile.html", data)
}

func (c *Controller) handlePostDetail(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, c.basePath+"/feed/post/")
	idStr = strings.Trim(idStr, "/")

	postID, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil || postID == 0 {
		http.NotFound(w, r)
		return
	}

	post, err := c.repo.GetPostByID(r.Context(), uint(postID))
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			http.NotFound(w, r)
			return
		}
		c.renderError(w, http.StatusInternalServerError, "could not load post")
		return
	}

	item := buildFeedItem(*post, c.basePath)
	data := PostPageData{
		Title:    "Post #" + uintToString(post.PostID),
		BasePath: c.basePath,
		Item:     item,
		Trending: c.loadTrending(r.Context()),
	}

	c.execute(w, "post.html", data)
}

func (c *Controller) buildItems(posts []Post) []FeedItem {
	items := make([]FeedItem, 0, len(posts))
	for _, p := range posts {
		items = append(items, buildFeedItem(p, c.basePath))
	}
	return items
}

func (c *Controller) loadTrending(ctx context.Context) []TrendingPond {
	tags, err := c.repo.ListTrendingPonds(ctx, 10)
	if err != nil {
		c.logger.Printf("list trending ponds: %v", err)
		return nil
	}

	for i := range tags {
		tags[i].URL = c.basePath + "/ponds/" + normalizeTagName(tags[i].TagName)
	}

	return tags
}

func (c *Controller) currentUser(r *http.Request) *account.IdentityUser {
	if c.authRepo == nil {
		return nil
	}

	user, ok := c.authRepo.CurrentUserFromRequest(r.Context(), r, "gs_auth_session")
	if !ok {
		return nil
	}

	return user
}

func (c *Controller) execute(w http.ResponseWriter, name string, data any) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if err := c.tmpl.ExecuteTemplate(w, name, data); err != nil {
		c.logger.Printf("template %s: %v", name, err)
		http.Error(w, "template execution failed", http.StatusInternalServerError)
	}
}

func (c *Controller) renderError(w http.ResponseWriter, status int, msg string) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(status)
	_ = c.tmpl.ExecuteTemplate(w, "error.html", map[string]any{
		"Status":  status,
		"Message": msg,
	})
}

func (c *Controller) handleCreatePost(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusMethodNotAllowed)
		json.NewEncoder(w).Encode(map[string]string{"error": "Method not allowed"})
		return
	}

	// Check authentication
	currentUser := c.currentUser(r)
	if currentUser == nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]string{"error": "Unauthorized"})
		return
	}

	// Parse form data
	if err := r.ParseForm(); err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Invalid form data"})
		return
	}

	content := strings.TrimSpace(r.FormValue("content"))
	if content == "" {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Content cannot be empty"})
		return
	}

	// Create the post
	post, err := c.repo.CreatePost(r.Context(), currentUser.UserID, content, "post")
	if err != nil {
		c.logger.Printf("failed to create post: %v", err)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "Failed to create post"})
		return
	}

	// Return success response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"message": "Post created successfully",
		"post_id": post.PostID,
	})
}