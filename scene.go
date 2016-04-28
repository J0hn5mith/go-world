package main

/*
Pointer non pointer mess
*/
type Scene struct {
	objects []*Object
    program uint32
}

func NewScene(program uint32) (*Scene) {
	scene := new(Scene)
    scene.program = program
	return scene

}

func (scene *Scene) addObject(object *Object) {
    scene.objects = append(scene.objects, object)
}
