package go_world

import (
	"github.com/go-gl/gl/v4.1-core/gl"
	"strings"
	"fmt"
)
func compileShader(source string, shaderType uint32) (uint32, error) {
	shader := gl.CreateShader(shaderType)

	csources, free := gl.Strs(source)
	gl.ShaderSource(shader, 1, csources, nil)
	free()
	gl.CompileShader(shader)

	var status int32
	gl.GetShaderiv(shader, gl.COMPILE_STATUS, &status)
	if status == gl.FALSE {
		var logLength int32
		gl.GetShaderiv(shader, gl.INFO_LOG_LENGTH, &logLength)

		log := strings.Repeat("\x00", int(logLength+1))
		gl.GetShaderInfoLog(shader, logLength, nil, gl.Str(log))

		return 0, fmt.Errorf("failed to compile %v: %v", source, log)
	}

	return shader, nil
}

var vertexShader = `
#version 330
in vec3 position;
uniform mat4 model;
uniform mat4 camera;
uniform mat4 projection;

void main() {
	gl_Position = projection * camera  * model * vec4(position, 1);
}
` + "\x00"

var fragmentShader = `
#version 330
uniform vec3 modelColor;

out vec4 outputColor;
void main() {
	outputColor = vec4(modelColor, 0);
}
` + "\x00"
