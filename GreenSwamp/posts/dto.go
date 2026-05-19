package posts

import (
	"crypto/rand"
	"encoding/hex"
	"html/template"
	"regexp"
	"strconv"
	"strings"

	account "github.com/venexene/serv-prog-go/greenswamp/account"
)

type NavLinks struct {
	Feed    string
	Profile string
	Post    string
}

type FeedItem struct {
	Post
	MediaURL           string
	MediaType          string
	ThumbnailURL       string
	AltText            string
	ContentHTML        template.HTML
	TimeLabel          string
	IsEvent            bool
	EventTimeLabel     string
	Nav                NavLinks
	LikedByCurrentUser bool `json:"liked_by_current_user"`
}

type FeedPage struct {
	Title       string
	BasePath    string
	Tag         string
	Items       []FeedItem
	Trending    []TrendingPond
	CurrentUser *account.IdentityUser
	CSRFToken   string
}

type ProfilePage struct {
	User      User
	Posts     []FeedItem
	PostCount int
	Avatar    string
	Bio       string
}

type ProfilePageData struct {
	Title       string
	BasePath    string
	Profile     ProfilePage
	Trending    []TrendingPond
	CurrentUser *account.IdentityUser
}

type PostPageData struct {
	Title    string
	BasePath string
	Item     FeedItem
	Trending []TrendingPond
	Comments []CommentItem
	CSRFToken string
}

type InteractionRequest struct {
	PostID uint   `json:"post_id"`
	Type   string `json:"type"`
}

type CommentRequest struct {
	PostID  uint   `json:"post_id"`
	Content string `json:"content"`
}

type CommentItem struct {
	Interaction
	AuthorName    string `json:"author_name"`
	AuthorAvatar  string `json:"author_avatar"`
	TimeLabel     string `json:"time_label"`
}

var hashtagRe = regexp.MustCompile(`(^|[^[:alnum:]_])#([[:alnum:]_]+)`)

func buildFeedItem(p Post, basePath string, userInteractions map[uint][]string) FeedItem {
	item := FeedItem{
		Post:         p,
		MediaURL:     stringOrEmpty(p.MediaURL),
		MediaType:    mediaKind(p.PostType, p.MediaType),
		ThumbnailURL: stringOrEmpty(p.ThumbnailURL),
		AltText:      stringOrEmpty(p.AltText),
		ContentHTML:  renderContentHTML(p.Content, basePath),
		TimeLabel:    formatTime(p.CreatedAt),
		IsEvent:      p.Event != nil,
		Nav: NavLinks{
			Feed:    basePath + "/feed",
			Profile: basePath + "/profile/" + p.User.Username,
			Post:    basePath + "/feed/post/" + uintToString(p.PostID),
		},
	}

	if p.Event != nil {
		item.EventTimeLabel = formatTime(p.Event.EventTime)
	}

	if types, ok := userInteractions[p.PostID]; ok {
		for _, t := range types {
			if t == "like" {
				item.LikedByCurrentUser = true
				break
			}
		}
	}

	return item
}

func renderContentHTML(content, basePath string) template.HTML {
	trimmed := strings.TrimSpace(content)
	if trimmed == "" {
		return ""
	}

	escaped := template.HTMLEscapeString(trimmed)
	matches := hashtagRe.FindAllStringSubmatchIndex(escaped, -1)
	if len(matches) == 0 {
		return template.HTML(escaped)
	}

	var out strings.Builder
	last := 0

	for _, idx := range matches {
		prefixEnd := idx[3]
		tagStart := idx[4]
		tagEnd := idx[5]
		tag := escaped[tagStart:tagEnd]

		out.WriteString(escaped[last:prefixEnd])

		href := basePath + "/ponds/" + strings.ToLower(tag)
		out.WriteString(`<a href="` + href + `" class="text-swamp-700 hover:text-swamp-500">#` + tag + `</a>`)

		last = idx[1]
	}

	out.WriteString(escaped[last:])
	return template.HTML(out.String())
}

func uintToString(v uint) string {
	return strconv.FormatUint(uint64(v), 10)
}

func stringOrEmpty(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}

func avatarOrFallback(url *string) string {
	if url != nil && strings.TrimSpace(*url) != "" {
		return *url
	}
	return "/static/avatars/avatar.jpeg"
}

func bioOrEmpty(bio *string) string {
	if bio == nil {
		return ""
	}
	return *bio
}

func mediaKind(postType string, mediaType *string) string {
	if mediaType != nil && strings.TrimSpace(*mediaType) != "" {
		return *mediaType
	}
	return postType
}

func normalizeTagName(tag string) string {
	tag = strings.TrimSpace(tag)
	tag = strings.TrimPrefix(tag, "#")
	return strings.ToLower(tag)
}

func randomString(n int) string {
	b := make([]byte, n)
	rand.Read(b)
	return hex.EncodeToString(b)[:n]
}

func formatTime(t interface{ Format(string) string }) string {
	return t.Format("02 Jan 2006 15:04")
}

func shortDate(t interface{ Format(string) string }) string {
	return t.Format("Jan 2")
}

func monthAbbr(t interface{ Format(string) string }) string {
	return t.Format("Jan")
}

func dayNum(t interface{ Format(string) string }) string {
	return t.Format("2")
}
