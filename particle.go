package main

import (
	"github.com/go-gl/mathgl/mgl32"
)

type Particle struct {
	object   *Object
	velocity mgl32.Vec3
	mass     float32
}

func NewParticle(scene *Scene) *Particle {
	particle := new(Particle)

	particle.object = NewObject(createCircleGeometry(60))
	particle.object.configure(scene.program)
	scene.addObject(particle.object)

	return particle
}

type ParticleSystem struct {
	particles []*Particle
	scene     *Scene
}

func NewParticleSystem(scene *Scene) *ParticleSystem {
	particleSystem := new(ParticleSystem)
	particleSystem.scene = scene

	return particleSystem
}

func (particleSystem *ParticleSystem) newParticle() *Particle {
	particle := NewParticle(particleSystem.scene)
	particleSystem.particles = append(particleSystem.particles, particle)

	return particle
}

func (particleSystem *ParticleSystem) animate(time_delta float32) {
	for _, particle := range particleSystem.particles {
		particle.object.position[0] += (float32)(particle.velocity[0] * time_delta)
		particle.object.position[1] += (float32)(particle.velocity[1] * time_delta)
		particle.object.position[2] += (float32)(particle.velocity[2] * time_delta)
	}
}
