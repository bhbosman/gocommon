package uiCommon

type UiStarted struct {
	Active bool
}

func NewUiStarted(active bool) *UiStarted {
	return &UiStarted{Active: active}
}

const UIState = "UI.State"
