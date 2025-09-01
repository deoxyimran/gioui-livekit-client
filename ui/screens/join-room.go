package screens

import (
	"bufio"
	"fmt"
	"image"
	"image/color"
	"log"
	"os/exec"
	"runtime"
	"strings"

	"gioui.org/f32"
	"gioui.org/io/event"
	"gioui.org/io/pointer"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"gioui.org/text"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
	"github.com/deoxyimran/gioui-livekit-client/ui/components"
	"github.com/deoxyimran/gioui-livekit-client/ui/mediasrcs/video"
	"github.com/deoxyimran/gioui-livekit-client/ui/res/icons"
	"github.com/deoxyimran/gioui-livekit-client/ui/state"
	mytheme "github.com/deoxyimran/gioui-livekit-client/ui/theme"
	mylayout "github.com/deoxyimran/gioui-livekit-client/ui/utils/layout"
	"github.com/deoxyimran/gioui-livekit-client/ui/utils/svg"
	"github.com/oligo/gioview/menu"
	"github.com/oligo/gioview/theme"
)

type vidDevice struct {
	name string
	desc string
}

func listVideoDevices() (devices []vidDevice) {
	var cmd *exec.Cmd
	// var vidFormat string
	switch runtime.GOOS {
	case "linux":
		// vidFormat = "v4l2"
		cmd = exec.Command("v4l2-ctl", "--list-devices") // needs v4l-utils
	case "windows":
		// vidFormat = "dshow"
		cmd = exec.Command("ffmpeg", "-hide_banner", "-list_devices", "true", "-f", "dshow", "-i", "dummy")
	}
	out, err := cmd.CombinedOutput()
	if err != nil {
		log.Println("Error listing video devices: ", err)
		return nil
	}
	scanner := bufio.NewScanner(strings.NewReader(string(out)))
	var taken bool
	var device vidDevice
	switch runtime.GOOS {
	case "linux":
		for scanner.Scan() {
			line := strings.TrimSpace(scanner.Text())
			if strings.HasSuffix(line, ":") {
				// new device block
				taken = false
				strs := strings.Split(line, ":")
				device.desc = strings.TrimSpace(strs[0])
			} else if strings.HasPrefix(line, "/dev/video") && !taken {
				taken = true // only first /dev/videoN per device
				device.name = strings.TrimSpace(line)
				devices = append(devices, device)
			}
		}
	case "windows":
	}
	return
}

type JoinRoom struct {
	// Widgets
	th               *material.Theme
	userNameEditor   editor
	deviceSetting    deviceSetting
	joinRoomClickble widget.Clickable
	// States/control vars
	stateManager *state.App
	vidCanvas    components.VideoCanvas
}

func NewJoinRoomScreen(stateManager *state.App) *JoinRoom {
	th := material.NewTheme()
	vs := video.NewWebcamSource("")
	j := &JoinRoom{
		stateManager:   stateManager,
		th:             th,
		vidCanvas:      components.NewVideoCanvas(&vs, image.Pt(380, 260)),
		userNameEditor: newEditor(th, "Enter a name", false),
	}
	j.deviceSetting = newDevSetting(th, stateManager, []string{"None"}, listVideoDevices())
	return j
}

func (j *JoinRoom) StopVideoCapture() {

}

type editor struct {
	th          *material.Theme
	placeholder string
	isPassword  bool
	edit        widget.Editor
}

func newEditor(th *material.Theme, placeholder string, isPassword bool) editor {
	return editor{th: th, placeholder: placeholder, isPassword: isPassword}
}

func (e *editor) layout(gtx C) D {
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

type iconButton struct {
	th      *material.Theme
	icon    image.Image
	onClick func()
}

func newIconButton(th *material.Theme, icon image.Image) iconButton {
	return iconButton{th: th, icon: icon}
}

func (i *iconButton) layout(gtx C) D {
	dims := layout.Background{}.Layout(gtx,
		// Background
		func(gtx C) D {
			defer clip.UniformRRect(image.Rectangle{Max: gtx.Constraints.Min}, 5).Push(gtx.Ops).Pop()
			event.Op(gtx.Ops, i)
			orig := color.NRGBA{30, 30, 30, 255}
			c := orig
			// Process pointer hover and click events
			for {
				ev, ok := gtx.Source.Event(pointer.Filter{
					Target: i,
					Kinds:  pointer.Enter | pointer.Leave | pointer.Press | pointer.Release,
				})
				if !ok {
					break
				}
				if x, ok := ev.(pointer.Event); ok {
					switch x.Kind {
					case pointer.Enter:
						c = color.NRGBA{15, 15, 15, 255} // lighter shade
					case pointer.Leave:
						c = orig // back to normal
					case pointer.Release:
						if i.onClick != nil {
							i.onClick()
						}
						gtx.Execute(op.InvalidateCmd{})
					}
				}
			}
			paint.ColorOp{Color: c}.Add(gtx.Ops)
			paint.PaintOp{}.Add(gtx.Ops)
			pointer.CursorPointer.Add(gtx.Ops)

			return D{Size: gtx.Constraints.Min}
		},
		// Icon
		func(gtx C) D {
			fmt.Println("Icon cmin: ", gtx.Constraints.Min)
			return layout.Inset{Left: unit.Dp(3), Right: unit.Dp(3)}.Layout(gtx, widget.Image{Src: paint.NewImageOp(i.icon), Position: layout.Center}.Layout)
		},
	)
	return D{Size: dims.Size}
}

type devToggleBtn struct {
	th              *material.Theme
	offIcon, onIcon image.Image
	text            string
	isActive        bool
	toggleFunc      func()
}

func newToggleButton(th *material.Theme, offIcon, onIcon image.Image, text string) devToggleBtn {
	return devToggleBtn{
		th:      th,
		offIcon: offIcon,
		onIcon:  onIcon,
		text:    text,
	}
}

func (tb *devToggleBtn) layout(gtx C) D {
	return layout.Background{}.Layout(gtx,
		// Background
		func(gtx C) D {
			defer clip.UniformRRect(image.Rectangle{Max: gtx.Constraints.Min}, 5).Push(gtx.Ops).Pop() // Rounded Rect
			event.Op(gtx.Ops, &tb.text)
			orig := color.NRGBA{30, 30, 30, 255}
			c := orig
			// Process pointer hover and click events
			for {
				ev, ok := gtx.Source.Event(pointer.Filter{
					Target: &tb.text,
					Kinds:  pointer.Enter | pointer.Leave | pointer.Press,
				})
				if !ok {
					break
				}
				if x, ok := ev.(pointer.Event); ok {
					switch x.Kind {
					case pointer.Enter:
						c = color.NRGBA{15, 15, 15, 255} // lighter shade
					case pointer.Leave:
						c = orig // back to normal
					case pointer.Press:
						tb.isActive = !tb.isActive
						if tb.toggleFunc != nil {
							tb.toggleFunc()
						}
						gtx.Execute(op.InvalidateCmd{})
					}
				}
			}
			paint.ColorOp{Color: c}.Add(gtx.Ops)
			paint.PaintOp{}.Add(gtx.Ops)
			pointer.CursorPointer.Add(gtx.Ops)

			return D{Size: gtx.Constraints.Min}
		},
		func(gtx C) D {
			if tb.isActive {
				return layout.UniformInset(unit.Dp(7)).Layout(gtx,
					func(gtx C) D {
						return layout.Flex{
							Axis:    layout.Horizontal,
							Spacing: layout.SpaceEvenly,
						}.Layout(gtx,
							layout.Rigid(widget.Image{Src: paint.NewImageOp(tb.onIcon), Position: layout.Center}.Layout),
							layout.Flexed(1, func(gtx C) D {
								lb := material.Label(tb.th, unit.Sp(14), tb.text)
								lb.Color = color.NRGBA{255, 255, 255, 255}
								lb.Alignment = text.Middle
								return layout.Center.Layout(gtx, lb.Layout)
							}),
						)
					},
				)
			} else {
				return layout.UniformInset(unit.Dp(7)).Layout(gtx,
					func(gtx C) D {
						return layout.Flex{
							Axis:    layout.Horizontal,
							Spacing: layout.SpaceEvenly,
						}.Layout(gtx,
							layout.Rigid(widget.Image{Src: paint.NewImageOp(tb.offIcon), Position: layout.Center}.Layout),
							layout.Flexed(1, func(gtx C) D {
								lb := material.Label(tb.th, unit.Sp(14), tb.text)
								lb.Color = color.NRGBA{255, 255, 255, 255}
								lb.Alignment = text.Middle
								return layout.Center.Layout(gtx, lb.Layout)
							}),
						)
					},
				)
			}
		},
	)
}

type deviceSetting struct {
	th *material.Theme
	st *state.App

	camToggleBtn   devToggleBtn
	camDropdownBtn iconButton
	camDropdownTh  *theme.Theme
	camDropdown    *menu.DropdownMenu

	micToggleBtn   devToggleBtn
	micDropdownBtn iconButton
	micDropdownTh  *theme.Theme
	micDropdown    *menu.DropdownMenu

	micPaths   []string
	vidDevices []vidDevice
}

func newDevSetting(th *material.Theme, st *state.App, micPaths []string, vidDevices []vidDevice) deviceSetting {
	d := deviceSetting{
		th: th,
		st: st,
	}
	sz := image.Pt(20, 20)
	camOffIcon, _ := svg.LoadSvg(strings.NewReader(icons.CamOff), sz)
	camOnIcon, _ := svg.LoadSvg(strings.NewReader(icons.CamOn), sz)
	micOnIcon, _ := svg.LoadSvg(strings.NewReader(icons.MicOn), sz)
	micOffIcon, _ := svg.LoadSvg(strings.NewReader(icons.MicOff), sz)
	angleDownIcon, _ := svg.LoadSvg(strings.NewReader(icons.AngleDown), sz)

	th_ := theme.NewTheme("./fonts", nil, false)
	th_.TextSize = unit.Sp(14)
	th_.Bg = color.NRGBA{255, 255, 255, 255}
	th_.Bg = color.NRGBA{140, 140, 140, 255}
	d.camDropdownTh = th_
	d.micDropdownTh = th_

	d.vidDevices = vidDevices
	d.micPaths = micPaths

	d.camDropdown = d.newCamDropdown()
	d.camToggleBtn = newToggleButton(th, camOffIcon, camOnIcon, "Camera")
	d.camDropdownBtn = newIconButton(th, angleDownIcon)

	d.micDropdown = d.newMicDropdown()
	d.micToggleBtn = newToggleButton(th, micOffIcon, micOnIcon, "Microphone")
	d.micDropdownBtn = newIconButton(th, angleDownIcon)

	return d
}

func (d *deviceSetting) newCamDropdown() *menu.DropdownMenu {
	options := [][]menu.MenuOption{}
	options = append(options, []menu.MenuOption{})
	if d.vidDevices == nil {
		d.vidDevices = []vidDevice{{name: "None", desc: "No camera found"}}
	}
	for i, v := range d.vidDevices {
		options[0] = append(options[0], menu.MenuOption{
			Layout: func(gtx menu.C, th *theme.Theme) menu.D {
				lb := material.Label(th.Theme, unit.Sp(14), v.desc)
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

func (d *deviceSetting) newMicDropdown() *menu.DropdownMenu {
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

func (d *deviceSetting) toggleMic() {
	d.st.MicOn = !d.st.MicOn
}

func (d *deviceSetting) toggleMicDropdown(gtx C) {
	d.micDropdown.ToggleVisibility(gtx)
}

func (d *deviceSetting) toggleCam() {
	d.st.CameraOn = !d.st.CameraOn
}

func (d *deviceSetting) toggleCamDropdown(gtx C) {
	d.camDropdown.ToggleVisibility(gtx)
}

func (d *deviceSetting) update(gtx C) {
	d.camDropdown.Update(gtx)
	d.micDropdown.Update(gtx)
}

func (d *deviceSetting) layout(gtx C, menuctx C) D {
	// Update states
	d.update(gtx)

	dims := layout.Flex{
		Axis:    layout.Horizontal,
		Spacing: layout.SpaceBetween,
	}.Layout(gtx,
		// Microphone
		layout.Flexed(1, func(gtx C) D {
			return layout.Flex{
				Axis: layout.Horizontal,
			}.Layout(gtx,
				// Toggle button
				layout.Flexed(1, func(gtx C) D {
					d.micToggleBtn.toggleFunc = d.toggleMic
					return d.micToggleBtn.layout(gtx)
				}),
				// Spacer
				layout.Rigid(layout.Spacer{Width: unit.Dp(3)}.Layout),
				// Dropdown btn and menu
				layout.Rigid(func(gtx C) D {
					d.micDropdownBtn.onClick = func() {
						d.toggleMicDropdown(gtx)
					}
					d.micDropdown.Layout(menuctx, d.micDropdownTh)
					return d.micDropdownBtn.layout(gtx)
				}),
				// Spacer
				layout.Rigid(layout.Spacer{Width: unit.Dp(10)}.Layout),
			)
		}),
		// Camera
		layout.Flexed(1, func(gtx C) D {
			return layout.Flex{
				Axis: layout.Horizontal,
			}.Layout(gtx,
				// Toggle button
				layout.Flexed(1, func(gtx C) D {
					d.camToggleBtn.toggleFunc = d.toggleCam
					return d.camToggleBtn.layout(gtx)
				}),
				// Spacer
				layout.Rigid(layout.Spacer{Width: unit.Dp(3)}.Layout),
				// Dropdown btn and menu
				layout.Rigid(func(gtx C) D {
					d.camDropdownBtn.onClick = func() {
						d.toggleCamDropdown(gtx)
					}
					d.camDropdown.Layout(menuctx, d.camDropdownTh)
					return d.camDropdownBtn.layout(gtx)
				}),
				// Spacer
				layout.Rigid(layout.Spacer{Width: unit.Dp(10)}.Layout),
			)
		}),
	)
	return dims
}

func (j *JoinRoom) Layout(gtx C, screenPointer *Screen) D {
	return layout.Background{}.Layout(gtx,
		// Fullscreen background
		func(gtx C) D {
			img := mytheme.AppBackground()
			defer clip.Rect{Max: gtx.Constraints.Max}.Push(gtx.Ops).Pop()
			scale := f32.Affine2D{}.Scale(f32.Point{}, f32.Point{
				X: float32(gtx.Constraints.Max.X) / float32(img.Bounds().Dx()),
				Y: float32(gtx.Constraints.Max.Y) / float32(img.Bounds().Dy()),
			})
			defer op.Affine(scale).Push(gtx.Ops).Pop()
			paint.NewImageOp(img).Add(gtx.Ops)
			paint.PaintOp{}.Add(gtx.Ops)
			return D{Size: gtx.Constraints.Max}
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
												// Video canvas layout
												layout.Rigid(j.vidCanvas.Layout),
												// Device selector layout
												layout.Rigid(func(gtx C) D {
													return layout.Inset{Top: unit.Dp(8)}.Layout(gtx, func(gtx C) D {
														menuctx := gtx
														menuctx.Constraints.Max.Y = 200
														fmt.Println("Menuctx maxcons: ", menuctx.Constraints.Max)
														c := gtx.Constraints
														c.Max.X = 500
														c.Min.Y, c.Max.Y = 32, 32
														gtx.Constraints = c
														return j.deviceSetting.layout(gtx, menuctx)
													})
												}),
												// Username editor layout
												layout.Rigid(func(gtx C) D {
													gtx.Constraints.Max.X = 300
													return j.userNameEditor.layout(gtx)
												}),
												// Take up remaining layout
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
