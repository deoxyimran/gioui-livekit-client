package screens

import (
	"fmt"
	"image"
	"image/color"
	"strings"

	"gioui.org/io/event"
	"gioui.org/io/pointer"
	"gioui.org/layout"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
	"github.com/deoxyimran/gioui-livekit-client/ui/components"
	"github.com/deoxyimran/gioui-livekit-client/ui/mediasrcs/video"
	"github.com/deoxyimran/gioui-livekit-client/ui/res/icons"
	mytheme "github.com/deoxyimran/gioui-livekit-client/ui/theme"
	mylayout "github.com/deoxyimran/gioui-livekit-client/ui/utils/layout"
	"github.com/deoxyimran/gioui-livekit-client/ui/utils/svg"
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
		deviceSelector: newDevSelector(th, stateManager, nil, nil),
		userNameEditor: newEditor(th, "Enter a name", false, 300),
	}

	return j
}

func (j *JoinRoom) StopVideoCapture() {

}

type editor struct {
	th          *material.Theme
	placeholder string
	isPassword  bool
	edit        widget.Editor
	w           int
}

func newEditor(th *material.Theme, placeholder string, isPassword bool, w int) editor {
	return editor{th: th, placeholder: placeholder, isPassword: isPassword, w: w}
}

func (e *editor) layout(gtx C) D {
	c := gtx.Constraints
	c.Max.X, c.Min.X = e.w, e.w
	gtx.Constraints = c
	edit := material.Editor(e.th, &e.edit, e.placeholder)
	edit.Color = color.NRGBA{R: 255, G: 255, B: 255, A: 255}
	edit.HintColor = color.NRGBA{R: 135, G: 135, B: 135, A: 220}
	edit.TextSize = unit.Sp(14)
	return layout.UniformInset(unit.Dp(10)).Layout(gtx,
		func(gtx C) D {
			return mylayout.Border(gtx, edit.Layout, 1, 8, color.NRGBA{R: 140, G: 140, B: 140, A: 255})
		},
	)
}

type toggleButton struct {
	th              *material.Theme
	offIcon, onIcon image.Image
	text            string
	isActive        bool
	onClick         func()
}

func newToggleButton(th *material.Theme, offIcon, onIcon image.Image, text string) toggleButton {
	return toggleButton{
		th:      th,
		offIcon: offIcon,
		onIcon:  onIcon,
		text:    text,
	}
}

func (tb *toggleButton) layout(gtx C) D {
	var dims D
	if tb.isActive {
		dims = layout.Background{}.Layout(gtx,
			// Background
			func(gtx C) D {
				defer clip.UniformRRect(image.Rectangle{Max: gtx.Constraints.Min}, 20).Push(gtx.Ops).Pop()
				color := color.NRGBA{30, 30, 30, 255}
				event.Op(gtx.Ops, &tb.text)
				// Check for pointer hover
				for {
					ev, ok := gtx.Source.Event(pointer.Filter{
						Target: &tb.text,
						Kinds:  pointer.Enter | pointer.Leave,
					})
					if !ok {
						break
					}
					if x, ok := ev.(pointer.Event); ok {
						switch x.Kind {
						case pointer.Enter:
							color = color.NRGBA{15, 15, 15, 255} // lighter shade
						case pointer.Leave:
							color = color.NRGBA{30, 30, 30, 255}
						}
					}
				}
				pointer.CursorPointer.Add(gtx.Ops)
				return D{Size: gtx.Constraints.Min}
			},
			func(gtx C) D {
				return D{}
			},
		)
	} else {
		// dims = material.IconButton(tb.th, &tb.clickable, tb.offIcon, tb.text).Layout(gtx)
	}
	return dims
}

type button struct {
	th   *material.Theme
	icon image.Image
}

type devSelector struct {
	th *material.Theme
	st *StateManager

	camToggleBtn   toggleButton
	camDropdownTh  *theme.Theme
	camDropdownBtn iconButton
	camDropdown    *menu.DropdownMenu

	micToggleBtn   toggleButton
	micDropdownTh  *theme.Theme
	micDropdownBtn iconButton
	micDropdown    *menu.DropdownMenu

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
	sz := image.Pt(32, 32)
	camOffIcon, _ := svg.LoadSvg(strings.NewReader(icons.CamOff), sz)
	camOnIcon, _ := svg.LoadSvg(strings.NewReader(icons.CamOn), sz)
	micOnIcon, _ := svg.LoadSvg(strings.NewReader(icons.MicOn), sz)
	micOffIcon, _ := svg.LoadSvg(strings.NewReader(icons.MicOff), sz)

	d.camDropdown = camDropDown
	// d.camDropdownBtn =
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
	for i, v := range d.micPaths {
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
			if d.camDropdownToggleBtn.Clicked(gtx) {
				d.camDropdown.ToggleVisibility(gtx)
			}
			return layout.Center.Layout(gtx, func(gtx C) D {
				dims := material.Button(d.th, &d.camDropdownToggleBtn, "Choose Cam").Layout(gtx)
				d.camDropdown.Layout(gtx, d.camDropdownTh)
				return dims
			})
		}),
		layout.Flexed(1, func(gtx C) D {
			d.micDropdown.Update(gtx)
			if d.micDropdownToggleBtn.Clicked(gtx) {
				d.micDropdown.ToggleVisibility(gtx)
			}
			return layout.Center.Layout(gtx, func(gtx C) D {
				dims := material.Button(d.th, &d.micDropdownToggleBtn, "Choose Mic").Layout(gtx)
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
					return layout.Flex{Axis: layout.Horizontal, Alignment: layout.Middle}.Layout(gtx,
						layout.Rigid(func(gtx C) D {
							return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
								layout.Rigid(func(gtx C) D {
									return layout.UniformInset(unit.Dp(10)).Layout(gtx,
										func(gtx C) D {
											return layout.Flex{
												Axis:      layout.Vertical,
												Alignment: layout.Middle,
											}.Layout(gtx,
												// Take up remaining
												layout.Flexed(1, func(gtx C) D {
													return D{Size: gtx.Constraints.Max}
												}),
												// Video canvas
												layout.Rigid(j.vidCanvas.Layout),
												// Device selector
												layout.Rigid(j.deviceSelector.layout),
												// // Username editor
												layout.Rigid(j.userNameEditor.layout),
												// Take up remaining
												layout.Flexed(1, func(gtx C) D {
													return D{Size: gtx.Constraints.Max}
												}),
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
