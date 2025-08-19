package screens

import (
	"fmt"
	"image"
	"image/color"

	"gioui.org/layout"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
	"github.com/deoxyimran/gioui-livekit-client/ui/components"
	"github.com/deoxyimran/gioui-livekit-client/ui/media/video"
	mytheme "github.com/deoxyimran/gioui-livekit-client/ui/theme"
	"github.com/deoxyimran/gioui-livekit-client/ui/utils"
	"github.com/oligo/gioview/menu"
	"github.com/oligo/gioview/theme"
)

type JoinRoom struct {
	// Widgets
	th               *material.Theme
	userNameEditor   editor
	deviceSelector   devSelector
	joinRoomClickble widget.Clickable
	// States/control vars
	stateManager *StateManager
	vidCanvas    components.VideoCanvas
}

func NewJoinRoomScreen(stateManager *StateManager) *JoinRoom {
	th := material.NewTheme()
	vs := video.NewWebcamSource("")
	j := &JoinRoom{
		stateManager:   stateManager,
		th:             th,
		vidCanvas:      components.NewVideoCanvas(&vs, image.Pt(380, 260)),
		deviceSelector: newDevSelector(th, stateManager, []string{"mic1", "mic2"}, []string{"cam1", "cam2"}),
		userNameEditor: newEditor(th),
	}

	return j
}

func (j *JoinRoom) StopVideoCapture() {

}

type editor struct {
	th         *material.Theme
	text       string
	isPassword string
	edit       widget.Editor
	w          int
}

func newEditor(th *material.Theme) editor {
	return editor{th: th}
}

func (e *editor) layout(gtx C) D {
	c := gtx.Constraints
	c.Max.X, c.Min.X = e.w, e.w
	gtx.Constraints = c
	edit := material.Editor(e.th, &e.edit, "Enter a name")
	edit.Color = color.NRGBA{R: 255, G: 255, B: 255, A: 255}
	edit.HintColor = color.NRGBA{R: 135, G: 135, B: 135, A: 220}
	edit.TextSize = unit.Sp(14)
	return layout.UniformInset(unit.Dp(10)).Layout(gtx,
		func(gtx C) D {
			return utils.BorderLayout(gtx, edit.Layout, 1, 8, color.NRGBA{R: 140, G: 140, B: 140, A: 255})
		},
	)
}

type devSelector struct {
	th *material.Theme
	st *StateManager

	camDropdownTh        *theme.Theme
	camDropdownClickable widget.Clickable
	camDropdown          *menu.DropdownMenu

	micDropdownTh        *theme.Theme
	micDropdownClickable widget.Clickable
	micDropdown          *menu.DropdownMenu

	micPaths []string
	camPaths []string
	w        int
}

func newDevSelector(th *material.Theme, st *StateManager, micPaths []string, camPaths []string) devSelector {
	d := devSelector{
		th: th,
	}
	var camDropDown, micDropDown *menu.DropdownMenu
	if camPaths != nil {
		d.camPaths = camPaths
		camDropDown = d.newCamDropdown()
	}
	if micPaths != nil {
		d.micPaths = micPaths
		micDropDown = d.newMicDropdown()
	}
	d.camDropdown = camDropDown
	d.micDropdown = micDropDown

	th_ := theme.NewTheme("./fonts", nil, false)
	th_.TextSize = unit.Sp(14)
	th_.Bg = color.NRGBA{255, 255, 255, 255}
	th_.Bg = color.NRGBA{140, 140, 140, 255}
	d.camDropdownTh = th_
	d.micDropdownTh = th_

	return d
}

func (d *devSelector) newCamDropdown() *menu.DropdownMenu {
	options := [][]menu.MenuOption{}
	options = append(options, []menu.MenuOption{})
	for i, v := range d.camPaths {
		options[0] = append(options[0], menu.MenuOption{
			Layout: func(gtx menu.C, th *theme.Theme) menu.D {
				lb := material.Label(th.Theme, unit.Sp(14), fmt.Sprint(v, i+1))
				return lb.Layout(gtx)
			},
			OnClicked: func() error {
				fmt.Println("Clicked cam", i+1)
				return nil
			},
		})
	}
	return menu.NewDropdownMenu(options)
}

func (d *devSelector) newMicDropdown() *menu.DropdownMenu {
	options := [][]menu.MenuOption{}
	options = append(options, []menu.MenuOption{})
	for i, v := range d.camPaths {
		options[0] = append(options[0], menu.MenuOption{
			Layout: func(gtx menu.C, th *theme.Theme) menu.D {
				lb := material.Label(th.Theme, unit.Sp(14), fmt.Sprint(v, i+1))
				return lb.Layout(gtx)
			},
			OnClicked: func() error {
				fmt.Println("Clicked mic", i+1)
				return nil
			},
		})
	}
	return menu.NewDropdownMenu(options)
}

func (d *devSelector) layout(gtx C) D {
	dims := layout.Flex{
		Axis: layout.Horizontal,
	}.Layout(gtx,
		layout.Flexed(1, func(gtx C) D {
			d.camDropdown.Update(gtx)
			if d.camDropdownClickable.Clicked(gtx) {
				d.camDropdown.ToggleVisibility(gtx)
			}
			return layout.Center.Layout(gtx, func(gtx C) D {
				dims := material.Button(d.th, &d.camDropdownClickable, "Choose Cam").Layout(gtx)
				d.camDropdown.Layout(gtx, d.camDropdownTh)
				return dims
			})
		}),
		layout.Flexed(1, func(gtx C) D {
			d.micDropdown.Update(gtx)
			if d.micDropdownClickable.Clicked(gtx) {
				d.micDropdown.ToggleVisibility(gtx)
			}
			return layout.Center.Layout(gtx, func(gtx C) D {
				dims := material.Button(d.th, &d.micDropdownClickable, "Choose Mic").Layout(gtx)
				d.micDropdown.Layout(gtx, d.micDropdownTh)
				return dims
			})
		}),
	)
	return dims
}

func (j *JoinRoom) Layout(gtx C, screenPointer *Screen) D {
	return layout.Background{}.Layout(gtx,
		// Fullscreen background
		func(gtx C) D {
			defer clip.Rect{Max: gtx.Constraints.Max}.Push(gtx.Ops).Pop()
			color := mytheme.BackgroundColor()
			paint.ColorOp{Color: color}.Add(gtx.Ops)
			paint.PaintOp{}.Add(gtx.Ops)
			return layout.Dimensions{Size: gtx.Constraints.Max}
		},
		// Main content
		func(gtx C) D {
			gtx.Constraints.Min = image.Pt(0, 0) // Reset Constraints Min
			return layout.UniformInset(10).Layout(gtx,
				func(gtx C) D {
					return layout.Flex{Axis: layout.Horizontal, Alignment: layout.Middle, Spacing: layout.SpaceAround}.Layout(gtx,
						layout.Rigid(func(gtx C) D {
							return layout.Flex{Axis: layout.Vertical, Alignment: layout.Middle}.Layout(gtx,
								layout.Rigid(func(gtx C) D {
									return layout.UniformInset(unit.Dp(10)).Layout(gtx,
										func(gtx C) D {
											return layout.Flex{
												Axis:      layout.Vertical,
												Alignment: layout.Middle,
											}.Layout(gtx,
												// Video canvas
												layout.Rigid(j.vidCanvas.Layout),
												// Device selector
												layout.Rigid(func(gtx C) D {
													// gtx.Constraints.Max.X, gtx.Constraints.Min.X = 600, 600
													// gtx.Constraints.Max.Y, gtx.Constraints.Min.Y = 200, 200
													return j.deviceSelector.layout(gtx)
												}),
												// // Username editor
												// layout.Rigid(j.userNameEditor.layout),
											)
										},
									)
								}),
							)
						}),
					)
				},
			)
		},
	)
}
