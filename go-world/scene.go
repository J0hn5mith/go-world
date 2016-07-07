package go_world

import ()

/*
Pointer non pointer mess
*/
type Scene struct {
	objects []*Object
	program uint32
}

func NewScene(program uint32) *Scene {
	scene := new(Scene)
	scene.program = program
	return scene

}

func (scene *Scene) addObject(object *Object) {
	scene.objects = append(scene.objects, object)
}

func (scene *Scene) AddObject(object *Object) {
	scene.objects = append(scene.objects, object)
}

func (scene *Scene) AddObjects(objects ...*Object) {
	for _, object := range objects {
        scene.AddObject(object)
    }
}

func (s *Scene) RemoveObject(object *Object) *Scene {
	for i, item := range s.objects {
		if item == object {
            first := s.objects[0:i]
            second := s.objects[i+1:]
            s.objects = append(first, second...)
	        return s
		}
	}
	return s
}

func (s *Scene) Objects() []*Object {
	return s.objects
}
