package go_world

import (
	"github.com/go-gl/glfw/v3.1/glfw"
)

/* Create a window for rendering a scene.  */
func CreateWindow(width int, height int, name string) (*glfw.Window, error) {

	glfw.WindowHint(glfw.Resizable, glfw.False)
	glfw.WindowHint(glfw.ContextVersionMajor, 4)
	glfw.WindowHint(glfw.ContextVersionMinor, 1)
	glfw.WindowHint(glfw.OpenGLProfile, glfw.OpenGLCoreProfile)
	glfw.WindowHint(glfw.OpenGLForwardCompatible, glfw.True)
	window, err := glfw.CreateWindow(width, height, name, nil, nil)

	if err != nil {
		panic(err)
	}

	window.MakeContextCurrent()
	return window, err
}
