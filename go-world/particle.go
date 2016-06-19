package go_world

import (
	"github.com/go-gl/mathgl/mgl32"
)

type Particle struct {
	object   *Object
	velocity mgl32.Vec3
	mass     float32
	radius   float32
	scene    *Scene
}

func NewParticle(scene *Scene) *Particle {
	particle := new(Particle)
	particle.radius = 1
	particle.scene = scene
	particle.createObject()

	return particle
}

func (p *Particle) SetVelocity(x, y, z float32) *Particle{
	p.velocity[0] = x
	p.velocity[1] = y
	p.velocity[2] = z
    return p
}

func (p *Particle) Velocity() mgl32.Vec3 {
	return p.velocity
}

func (p *Particle) ApplyForce(x, y, z float32) {
	p.velocity[0] = p.velocity[0] + x
	p.velocity[1] = p.velocity[1] + y
	p.velocity[2] = p.velocity[2] + z
}

func (p Particle) Object() *Object {
	return p.object
}

func (p Particle) Radius() float32 {
	return p.radius
}

func (p *Particle) SetRadius(radius float32) *Particle {
	p.radius = radius
	p.createObject()
	return p
}

func (p *Particle) Position() mgl32.Vec3 {
	return p.object.Position()
}

func (p *Particle) SetPosition(x, y, z float32) *Particle{
	p.object.SetPosition(x, y, z)
    return p
}

func (p *Particle) Mass() float32 {
	return 1
}

func (p *Particle) createObject() {
	//if p.object != nil {
		//p.scene.RemoveObject(p.object)
	//}
	//p.object = NewObject(
        //createCircleGeometry(60, p.radius).Load(p.scene.program),
    //)
	//p.scene.addObject(p.object)
}

type ParticleSystem struct {
	particles          []*Particle
	forceFields        []ParticleForceField
	constraints        []Constraint
	collisionHandler   CollisionHandler
	gravitationHandler GravitationHandler
	scene              *Scene
}

func NewParticleSystem(scene *Scene) *ParticleSystem {
	particleSystem := new(ParticleSystem)
	particleSystem.scene = scene

	return particleSystem
}

func (particleSystem *ParticleSystem) NewParticle() *Particle {
	particle := NewParticle(particleSystem.scene)
	particleSystem.particles = append(particleSystem.particles, particle)

	return particle
}

func (ps *ParticleSystem) RemoveParticle(particle *Particle) {
	for i, item := range ps.particles {
		if item == particle {
			first := ps.particles[0:i]
			second := ps.particles[i+1:]
			ps.scene.RemoveObject(item.object)
			ps.particles = append(first, second...)
			return
		}
	}
}

func (ps *ParticleSystem) Particles() []*Particle {
	return ps.particles
}

func (ps *ParticleSystem) AddForceField(ff ParticleForceField) {
	ps.forceFields = append(ps.forceFields, ff)
}

func (ps *ParticleSystem) AddConstraint(c Constraint) {
	ps.constraints = append(ps.constraints, c)
}

func (ps *ParticleSystem) SetCollisionHandler(ch CollisionHandler) {
	ps.collisionHandler = ch
}

func (ps *ParticleSystem) SetGravitationHandler(gh GravitationHandler) {
	ps.gravitationHandler = gh
}

func (ps *ParticleSystem) Update(time_delta float32) {
	ps.applyForces(time_delta)
	ps.applyGravitation(time_delta)
	ps.animate(time_delta)
	ps.handleCollisions()
	ps.applyConstraints(time_delta)
}

func (ps *ParticleSystem) applyForces(time_delta float32) {
	for _, p := range ps.particles {
		for _, ff := range ps.forceFields {
			ff.Apply(p, time_delta)
		}
	}
}

func (ps *ParticleSystem) applyGravitation(time_delta float32) {
	if ps.gravitationHandler != nil {
		ps.gravitationHandler.Apply(ps.particles, time_delta)
	}
}

func (ps *ParticleSystem) applyConstraints(time_delta float32) {
	for _, p := range ps.particles {
		for _, c := range ps.constraints {
			c.Apply(p)
		}
	}
}

func (particleSystem *ParticleSystem) handleCollisions() {
	if particleSystem.collisionHandler != nil {
		particleSystem.collisionHandler.Apply(particleSystem.particles)
	}
}

func (particleSystem *ParticleSystem) animate(time_delta float32) {
	for _, particle := range particleSystem.particles {

		particle.object.position[0] += (float32)(particle.velocity[0] * time_delta)
		particle.object.position[1] += (float32)(particle.velocity[1] * time_delta)
		particle.object.position[2] += (float32)(particle.velocity[2] * time_delta)
	}
}

type ParticleForceField interface {
	Apply(p *Particle, time_detla float32)
}

type Constraint interface {
	Apply(p *Particle)
}

type CollisionHandler interface {
	Apply(p []*Particle)
}

type GravitationHandler interface {
	Apply(p []*Particle, time_delta float32)
}
