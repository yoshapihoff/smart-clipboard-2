package fyne_app

import (
	"fmt"
	"smart-clipboard-2/internal/config"
	"time"
)

type ClipboardItem struct {
	Content    string
	Timestamp  time.Time
	Preview    string
	ClickCount int
}

func (ci ClipboardItem) String() string {
	return fmt.Sprintf("%s (%d clicks)", ci.Preview, ci.ClickCount)
}

type ClipboardManager struct {
	items []ClipboardItem
	cfg   *config.Config
}

func (cm ClipboardManager) String() string {
	return fmt.Sprintf("%v", cm.items)
}

func NewClipboardManager() *ClipboardManager {
	cfg := config.GetConfig()
	return &ClipboardManager{
		items: make([]ClipboardItem, 0, cfg.ClipboardHistorySize),
		cfg:   cfg,
	}
}

func (cm *ClipboardManager) AddItem(content string) {
	found := false
	for i, item := range cm.items {
		if item.Content == content {
			cm.items[i].ClickCount = cm.items[i].ClickCount + 1
			found = true
			break
		}
	}
	if !found {
		cm.items = append(cm.items, ClipboardItem{
			Content:    content,
			Timestamp:  time.Now(),
			Preview:    getPreview(content),
			ClickCount: 0,
		})
	}
	cm.sortItems()
	if len(cm.items) > cm.cfg.ClipboardHistorySize {
		cm.items = cm.items[:cm.cfg.ClipboardHistorySize]
	}
}

func (cm *ClipboardManager) sortItems() {
	for i := 0; i < len(cm.items)-1; i++ {
		for j := i + 1; j < len(cm.items); j++ {
			if cm.shouldSwap(cm.items[i], cm.items[j]) {
				cm.items[i], cm.items[j] = cm.items[j], cm.items[i]
			}
		}
	}
}

func (cm *ClipboardManager) shouldSwap(a, b ClipboardItem) bool {
	if a.ClickCount != b.ClickCount {
		return a.ClickCount < b.ClickCount
	}
	return a.Timestamp.Before(b.Timestamp)
}

func getPreview(content string) string {
	if len(content) <= 32 {
		return content
	}
	return content[:32] + "..."
}
