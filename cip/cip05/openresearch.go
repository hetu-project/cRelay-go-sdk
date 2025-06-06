package cip04

import (
	"fmt"
	"strings"

	"github.com/nbd-wtf/go-nostr"
	"github.com/nbd-wtf/go-nostr/cip"
)

// PaperEvent represents a paper operation in openresearch subspace
type PaperEvent struct {
	*nostr.SubspaceOpEvent
	DOI       string
	PaperType string
	Authors   []string
	Keywords  []string
	Year      string
	Journal   string
}

// SetPaperInfo sets the paper information
func (e *PaperEvent) SetPaperInfo(doi, paperType string, authors, keywords []string, year, journal string) {
	e.DOI = doi
	e.PaperType = paperType
	e.Authors = authors
	e.Keywords = keywords
	e.Year = year
	e.Journal = journal

	e.Tags = append(e.Tags,
		nostr.Tag{"doi", doi},
		nostr.Tag{"paper_type", paperType},
		nostr.Tag{"year", year},
		nostr.Tag{"journal", journal},
	)

	if len(authors) > 0 {
		authorsTag := nostr.Tag{"authors"}
		authorsTag = append(authorsTag, authors...)
		e.Tags = append(e.Tags, authorsTag)
	}

	if len(keywords) > 0 {
		keywordsTag := nostr.Tag{"keywords"}
		keywordsTag = append(keywordsTag, keywords...)
		e.Tags = append(e.Tags, keywordsTag)
	}
}

// AnnotationEvent represents an annotation operation in openresearch subspace
type AnnotationEvent struct {
	*nostr.SubspaceOpEvent
	PaperID  string
	Position string
	Type     string
	ParentID string
}

// SetAnnotationInfo sets the annotation information
func (e *AnnotationEvent) SetAnnotationInfo(paperID, position, annotationType, parentID string) {
	e.PaperID = paperID
	e.Position = position
	e.Type = annotationType
	e.ParentID = parentID

	e.Tags = append(e.Tags,
		nostr.Tag{"paper_id", paperID},
		nostr.Tag{"position", position},
		nostr.Tag{"type", annotationType},
	)

	if parentID != "" {
		e.Tags = append(e.Tags, nostr.Tag{"parent", parentID})
	}
}

// ReviewEvent represents a review operation in openresearch subspace
type ReviewEvent struct {
	*nostr.SubspaceOpEvent
	PaperID string
	Rating  string
	Aspects map[string]string
	Content string
}

// SetReviewInfo sets the review information
func (e *ReviewEvent) SetReviewInfo(paperID, rating string, aspects map[string]string) {
	e.PaperID = paperID
	e.Rating = rating
	e.Aspects = aspects

	e.Tags = append(e.Tags,
		nostr.Tag{"paper_id", paperID},
		nostr.Tag{"rating", rating},
	)

	if len(aspects) > 0 {
		aspectsStr := ""
		for k, v := range aspects {
			if aspectsStr != "" {
				aspectsStr += ","
			}
			aspectsStr += fmt.Sprintf("%s:%s", k, v)
		}
		e.Tags = append(e.Tags, nostr.Tag{"aspects", aspectsStr})
	}
}

// AIAnalysisEvent represents an AI analysis operation in openresearch subspace
type AIAnalysisEvent struct {
	*nostr.SubspaceOpEvent
	AnalysisType string
	PaperIDs     []string
	Prompt       string
	Content      string
}

// SetAIAnalysisInfo sets the AI analysis information
func (e *AIAnalysisEvent) SetAIAnalysisInfo(analysisType string, paperIDs []string, prompt string) {
	e.AnalysisType = analysisType
	e.PaperIDs = paperIDs
	e.Prompt = prompt

	e.Tags = append(e.Tags,
		nostr.Tag{"analysis_type", analysisType},
		nostr.Tag{"prompt", prompt},
	)

	if len(paperIDs) > 0 {
		paperIDsTag := nostr.Tag{"paper_ids"}
		paperIDsTag = append(paperIDsTag, paperIDs...)
		e.Tags = append(e.Tags, paperIDsTag)
	}
}

// DiscussionEvent represents a discussion operation in openresearch subspace
type DiscussionEvent struct {
	*nostr.SubspaceOpEvent
	Topic      string
	ParentID   string
	References []string
	Content    string
}

// SetDiscussionInfo sets the discussion information
func (e *DiscussionEvent) SetDiscussionInfo(topic, parentID string, references []string) {
	e.Topic = topic
	e.ParentID = parentID
	e.References = references

	e.Tags = append(e.Tags,
		nostr.Tag{"topic", topic},
	)

	if parentID != "" {
		e.Tags = append(e.Tags, nostr.Tag{"parent", parentID})
	}

	if len(references) > 0 {
		referencesTag := nostr.Tag{"references"}
		referencesTag = append(referencesTag, references...)
		e.Tags = append(e.Tags, referencesTag)
	}
}

// ReadPaperEvent represents a paper reading operation in openresearch subspace
type ReadPaperEvent struct {
	*nostr.SubspaceOpEvent
	PaperID  string
	UserID   string
	Duration string
	Depth    string
}

// SetReadPaperInfo sets the read paper information
func (e *ReadPaperEvent) SetReadPaperInfo(paperID, userID, duration, depth string) {
	e.PaperID = paperID
	e.UserID = userID
	e.Duration = duration
	e.Depth = depth

	e.Tags = append(e.Tags,
		nostr.Tag{"paper_id", paperID},
		nostr.Tag{"user_id", userID},
		nostr.Tag{"duration", duration},
		nostr.Tag{"depth", depth},
	)
}

// CoCreatePaperEvent represents a collaborative paper creation operation in openresearch subspace
type CoCreatePaperEvent struct {
	*nostr.SubspaceOpEvent
	PaperID string
	UserIDs []string
	Quality string
	Content string
}

// SetCoCreatePaperInfo sets the collaborative paper creation information
func (e *CoCreatePaperEvent) SetCoCreatePaperInfo(paperID string, userIDs []string, quality string) {
	e.PaperID = paperID
	e.UserIDs = userIDs
	e.Quality = quality

	e.Tags = append(e.Tags,
		nostr.Tag{"paper_id", paperID},
		nostr.Tag{"quality", quality},
	)

	if len(userIDs) > 0 {
		userIDsTag := nostr.Tag{"user_ids"}
		userIDsTag = append(userIDsTag, userIDs...)
		e.Tags = append(e.Tags, userIDsTag)
	}
}

// ParseOpenResearchEvent parses a Nostr event into an openresearch event
func ParseOpenResearchEvent(evt nostr.Event) (nostr.SubspaceOpEventPtr, error) {
	// Extract common fields
	subspaceID := ""
	var authTag cip.AuthTag
	parents := []string{}

	for _, tag := range evt.Tags {
		if len(tag) < 2 {
			continue
		}
		switch tag[0] {
		case "sid":
			subspaceID = tag[1]
		case "auth":
			auth, err := cip.ParseAuthTag(tag[1])
			if err != nil {
				return nil, fmt.Errorf("failed to parse auth tag: %v", err)
			}
			authTag = auth
		case "parent":
			parents = append(parents, tag[1:]...)
		}
	}

	// Get operation from kind
	operation, exists := cip.GetOpFromKind(evt.Kind)
	if !exists {
		return nil, fmt.Errorf("unknown kind value: %d", evt.Kind)
	}

	// Parse based on operation type
	switch operation {
	case cip.OpPaper:
		return parsePaperEvent(evt, subspaceID, operation, authTag, parents)
	case cip.OpAnnotation:
		return parseAnnotationEvent(evt, subspaceID, operation, authTag, parents)
	case cip.OpReview:
		return parseReviewEvent(evt, subspaceID, operation, authTag, parents)
	case cip.OpAIAnalysis:
		return parseAIAnalysisEvent(evt, subspaceID, operation, authTag, parents)
	case cip.OpDiscussion:
		return parseDiscussionEvent(evt, subspaceID, operation, authTag, parents)
	case cip.OpReadPaper:
		return parseReadPaperEvent(evt, subspaceID, operation, authTag, parents)
	case cip.OpCoCreate:
		return parseCoCreatePaperEvent(evt, subspaceID, operation, authTag, parents)
	default:
		return nil, fmt.Errorf("unknown operation type: %s", operation)
	}
}

func parsePaperEvent(evt nostr.Event, subspaceID, operation string, authTag cip.AuthTag, parents []string) (*PaperEvent, error) {
	paper := &PaperEvent{
		SubspaceOpEvent: &nostr.SubspaceOpEvent{
			SubspaceID: subspaceID,
			Operation:  operation,
			AuthTag:    authTag,
			Event:      evt,
			Parents:    parents,
		},
	}
	paper.Event.Content = evt.Content

	for _, tag := range evt.Tags {
		if len(tag) < 2 {
			continue
		}
		switch tag[0] {
		case "doi":
			paper.DOI = tag[1]
		case "paper_type":
			paper.PaperType = tag[1]
		case "authors":
			paper.Authors = tag[1:]
		case "keywords":
			paper.Keywords = tag[1:]
		case "year":
			paper.Year = tag[1]
		case "journal":
			paper.Journal = tag[1]
		}
	}

	return paper, nil
}

func parseAnnotationEvent(evt nostr.Event, subspaceID, operation string, authTag cip.AuthTag, parents []string) (*AnnotationEvent, error) {
	annotation := &AnnotationEvent{
		SubspaceOpEvent: &nostr.SubspaceOpEvent{
			SubspaceID: subspaceID,
			Operation:  operation,
			AuthTag:    authTag,
			Event:      evt,
			Parents:    parents,
		},
	}
	annotation.Event.Content = evt.Content

	for _, tag := range evt.Tags {
		if len(tag) < 2 {
			continue
		}
		switch tag[0] {
		case "paper_id":
			annotation.PaperID = tag[1]
		case "position":
			annotation.Position = tag[1]
		case "type":
			annotation.Type = tag[1]
		case "parent":
			annotation.ParentID = tag[1]
		}
	}

	return annotation, nil
}

func parseReviewEvent(evt nostr.Event, subspaceID, operation string, authTag cip.AuthTag, parents []string) (*ReviewEvent, error) {
	review := &ReviewEvent{
		SubspaceOpEvent: &nostr.SubspaceOpEvent{
			SubspaceID: subspaceID,
			Operation:  operation,
			AuthTag:    authTag,
			Event:      evt,
			Parents:    parents,
		},
		Content: evt.Content,
		Aspects: make(map[string]string),
	}

	for _, tag := range evt.Tags {
		if len(tag) < 2 {
			continue
		}
		switch tag[0] {
		case "paper_id":
			review.PaperID = tag[1]
		case "rating":
			review.Rating = tag[1]
		case "aspects":
			aspects := tag[1]
			// Parse aspects string like "methodology:4,novelty:5"
			aspectPairs := strings.Split(aspects, ",")
			for _, pair := range aspectPairs {
				parts := strings.Split(pair, ":")
				if len(parts) == 2 {
					review.Aspects[parts[0]] = parts[1]
				}
			}
		}
	}

	return review, nil
}

func parseAIAnalysisEvent(evt nostr.Event, subspaceID, operation string, authTag cip.AuthTag, parents []string) (*AIAnalysisEvent, error) {
	analysis := &AIAnalysisEvent{
		SubspaceOpEvent: &nostr.SubspaceOpEvent{
			SubspaceID: subspaceID,
			Operation:  operation,
			AuthTag:    authTag,
			Event:      evt,
			Parents:    parents,
		},
		Content: evt.Content,
	}

	for _, tag := range evt.Tags {
		if len(tag) < 2 {
			continue
		}
		switch tag[0] {
		case "analysis_type":
			analysis.AnalysisType = tag[1]
		case "paper_ids":
			analysis.PaperIDs = tag[1:]
		case "prompt":
			analysis.Prompt = tag[1]
		}
	}

	return analysis, nil
}

func parseDiscussionEvent(evt nostr.Event, subspaceID, operation string, authTag cip.AuthTag, parents []string) (*DiscussionEvent, error) {
	discussion := &DiscussionEvent{
		SubspaceOpEvent: &nostr.SubspaceOpEvent{
			SubspaceID: subspaceID,
			Operation:  operation,
			AuthTag:    authTag,
			Event:      evt,
			Parents:    parents,
		},
		Content: evt.Content,
	}

	for _, tag := range evt.Tags {
		if len(tag) < 2 {
			continue
		}
		switch tag[0] {
		case "topic":
			discussion.Topic = tag[1]
		case "parent":
			discussion.ParentID = tag[1]
		case "references":
			discussion.References = tag[1:]
		}
	}

	return discussion, nil
}

func parseReadPaperEvent(evt nostr.Event, subspaceID, operation string, authTag cip.AuthTag, parents []string) (*ReadPaperEvent, error) {
	readPaper := &ReadPaperEvent{
		SubspaceOpEvent: &nostr.SubspaceOpEvent{
			SubspaceID: subspaceID,
			Operation:  operation,
			AuthTag:    authTag,
			Event:      evt,
			Parents:    parents,
		},
	}
	readPaper.Event.Content = evt.Content
	
	for _, tag := range evt.Tags {
		if len(tag) < 2 {
			continue
		}
		switch tag[0] {
		case "paper_id":
			readPaper.PaperID = tag[1]
		case "user_id":
			readPaper.UserID = tag[1]
		case "duration":
			readPaper.Duration = tag[1]
		case "depth":
			readPaper.Depth = tag[1]
		}
	}

	return readPaper, nil
}

func parseCoCreatePaperEvent(evt nostr.Event, subspaceID, operation string, authTag cip.AuthTag, parents []string) (*CoCreatePaperEvent, error) {
	coCreate := &CoCreatePaperEvent{
		SubspaceOpEvent: &nostr.SubspaceOpEvent{
			SubspaceID: subspaceID,
			Operation:  operation,
			AuthTag:    authTag,
			Event:      evt,
			Parents:    parents,
		},
		Content: evt.Content,
	}

	for _, tag := range evt.Tags {
		if len(tag) < 2 {
			continue
		}
		switch tag[0] {
		case "paper_id":
			coCreate.PaperID = tag[1]
		case "user_ids":
			coCreate.UserIDs = tag[1:]
		case "quality":
			coCreate.Quality = tag[1]
		}
	}

	return coCreate, nil
}

// NewPaperEvent creates a new paper event
func NewPaperEvent(subspaceID string) (*PaperEvent, error) {
	baseEvent, err := nostr.NewSubspaceOpEvent(subspaceID, cip.KindOpenResearchPaper)
	if err != nil {
		return nil, err
	}
	return &PaperEvent{
		SubspaceOpEvent: baseEvent,
	}, nil
}

// NewAnnotationEvent creates a new annotation event
func NewAnnotationEvent(subspaceID string) (*AnnotationEvent, error) {
	baseEvent, err := nostr.NewSubspaceOpEvent(subspaceID, cip.KindOpenResearchAnnotation)
	if err != nil {
		return nil, err
	}
	return &AnnotationEvent{
		SubspaceOpEvent: baseEvent,
	}, nil
}

// NewReviewEvent creates a new review event
func NewReviewEvent(subspaceID string) (*ReviewEvent, error) {
	baseEvent, err := nostr.NewSubspaceOpEvent(subspaceID, cip.KindOpenResearchReview)
	if err != nil {
		return nil, err
	}
	return &ReviewEvent{
		SubspaceOpEvent: baseEvent,
		Aspects:         make(map[string]string),
	}, nil
}

// NewAIAnalysisEvent creates a new AI analysis event
func NewAIAnalysisEvent(subspaceID string) (*AIAnalysisEvent, error) {
	baseEvent, err := nostr.NewSubspaceOpEvent(subspaceID, cip.KindOpenResearchAIAnalysis)
	if err != nil {
		return nil, err
	}
	return &AIAnalysisEvent{
		SubspaceOpEvent: baseEvent,
	}, nil
}

// NewDiscussionEvent creates a new discussion event
func NewDiscussionEvent(subspaceID string) (*DiscussionEvent, error) {
	baseEvent, err := nostr.NewSubspaceOpEvent(subspaceID, cip.KindOpenResearchDiscussion)
	if err != nil {
		return nil, err
	}
	return &DiscussionEvent{
		SubspaceOpEvent: baseEvent,
	}, nil
}

// NewReadPaperEvent creates a new read paper event
func NewReadPaperEvent(subspaceID string) (*ReadPaperEvent, error) {
	baseEvent, err := nostr.NewSubspaceOpEvent(subspaceID, cip.KindOpenResearchReadPaper)
	if err != nil {
		return nil, err
	}
	return &ReadPaperEvent{
		SubspaceOpEvent: baseEvent,
	}, nil
}

// NewCoCreatePaperEvent creates a new collaborative paper creation event
func NewCoCreatePaperEvent(subspaceID string) (*CoCreatePaperEvent, error) {
	baseEvent, err := nostr.NewSubspaceOpEvent(subspaceID, cip.KindOpenResearchCoCreate)
	if err != nil {
		return nil, err
	}
	return &CoCreatePaperEvent{
		SubspaceOpEvent: baseEvent,
	}, nil
}
