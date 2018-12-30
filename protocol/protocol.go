package protocol

import (
	"syscall"
)

const (
	Center Align = "center"
	Right  Align = "right"
	Left   Align = "left"

	SIGRT_MIN = syscall.Signal(34)
	SIGRT_MAX = syscall.Signal(65)
)

var (
	RealTimeSignals = map[string]syscall.Signal{
		"RTMIN":    SIGRT_MIN,
		"RTMIN+1":  SIGRT_MIN + 1,
		"RTMIN+2":  SIGRT_MIN + 2,
		"RTMIN+3":  SIGRT_MIN + 3,
		"RTMIN+4":  SIGRT_MIN + 4,
		"RTMIN+5":  SIGRT_MIN + 5,
		"RTMIN+6":  SIGRT_MIN + 6,
		"RTMIN+7":  SIGRT_MIN + 7,
		"RTMIN+8":  SIGRT_MIN + 8,
		"RTMIN+9":  SIGRT_MIN + 9,
		"RTMIN+10": SIGRT_MIN + 10,
		"RTMIN+11": SIGRT_MIN + 11,
		"RTMIN+12": SIGRT_MIN + 12,
		"RTMIN+13": SIGRT_MIN + 13,
		"RTMIN+14": SIGRT_MIN + 14,
		"RTMIN+15": SIGRT_MIN + 15,
		"RTMAX-14": SIGRT_MAX - 14,
		"RTMAX-13": SIGRT_MAX - 13,
		"RTMAX-12": SIGRT_MAX - 12,
		"RTMAX-11": SIGRT_MAX - 11,
		"RTMAX-10": SIGRT_MAX - 10,
		"RTMAX-9":  SIGRT_MAX - 9,
		"RTMAX-8":  SIGRT_MAX - 8,
		"RTMAX-7":  SIGRT_MAX - 7,
		"RTMAX-6":  SIGRT_MAX - 6,
		"RTMAX-5":  SIGRT_MAX - 5,
		"RTMAX-4":  SIGRT_MAX - 4,
		"RTMAX-3":  SIGRT_MAX - 3,
		"RTMAX-2":  SIGRT_MAX - 2,
		"RTMAX-1":  SIGRT_MAX - 1,
		"RTMAX":    SIGRT_MAX,
	}
)

type Header struct {
	Version        int  `json:"version"`
	StopSignal     int  `json:"stop_signal,omitempty"`
	ContinueSignal int  `json:"cont_signal,omitempty"`
	ClickEvents    bool `json:"click_events,omitempty"`
}

type Align string

type Block struct {
	FullText string `json:"full_text"`
	//ShortText           string `json:"short_text"`
	Color               string `json:"color"`
	Background          string `json:"background,omitempty"`
	Border              string `json:"border,omitempty"`
	MinWidth            int    `json:"min_width,omitempty"`
	Align               Align  `json:"align,omitempty"`
	Urgent              bool   `json:"urgent,omitempty"`
	Name                string `json:"name,omitempty"`
	Instance            string `json:"instance,omitempty"`
	Separator           bool   `json:"separator"`
	SeparatorBlockWidth int    `json:"separator_block_width,omitempty"`
}

type ClickEvent struct {
	Name      string   `json:"name"`
	Instance  string   `json:"instance"`
	Button    int      `json:"button"`
	Modifiers []string `json:"modifiers"`
	X         int      `json:"x"`
	Y         int      `json:"y"`
	RelativeX int      `json:"relative_x"`
	RelativeY int      `json:"relative_y"`
	Width     int      `json:"width"`
	Height    int      `json:"height"`
}
