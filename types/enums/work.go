package enums

type WorkType string

const (
	WorkTypeCopyright    WorkType = "copyright"
	WorkTypeRelatedRight WorkType = "relatedright"
	WorkTypeIP           WorkType = "ip"
)

func (w WorkType) String() string {
	return string(w)
}

func (w WorkType) Text() string {
	switch w {
	case WorkTypeCopyright:
		return "Quyền tác giả"
	case WorkTypeRelatedRight:
		return "Quuyền liên quan"
	case WorkTypeIP:
		return "Sỡ hữu trí tuệ"
	default:
		return ""
	}
}
