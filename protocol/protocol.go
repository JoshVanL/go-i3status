package protocol

const (
	Center Align = "center"
	Right  Align = "right"
	Left   Align = "left"
)

type Header struct {
	Version        int  `json:"version"`
	StopSignal     int  `json:"stop_signal,omitempty"`
	ContinueSignal int  `json:"cont_signal,omitempty"`
	ClickEvents    bool `json:"click_events,omitempty"`
}

type Align string

type Block struct {
	FullText            string `json:"full_text"`
	ShortText           string `json:"short_text"`
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
