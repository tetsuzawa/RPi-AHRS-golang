package main

import (
	"math"

	"github.com/westphae/quaternion"
)

func (p Parameters) updateAttitude(q *quaternion.Quaternion) {
	var (
		norm,               // norm to normalise
		bx, bz,             // reference direction of flux in earth frame
		qew, qex, qey, qez, // q_hat_dot_epsilon
		R11, R12, R13,      // rotation matrix
		R21, R22, R23,      //
		R31, R32, R33,      //
		f1, f2, f3,         // objective function
		a1, a2, a3, a4, a5, a6, a7, a8 float64 // auxiliary variables
	)

	// normalise the accelerometer measurement
	norm = math.Sqrt(p.ax*p.ax + p.ay*p.ay + p.az*p.az)
	p.ax /= norm
	p.ay /= norm
	p.az /= norm

	// normalise the magnetometer measurement
	norm = math.Sqrt(p.mx*p.mx + p.my*p.my + p.mz*p.mz)
	p.mx /= norm
	p.my /= norm
	p.mz /= norm

	for i := 0; i < 3000; i++ {
		// compute rotation matrix
		a1 = q.W*q.W - 0.5
		R11 = a1 + q.X*q.X
		R22 = a1 + q.Y*q.Y
		R33 = a1 + q.Z*q.Z
		a1 = q.X * q.Y
		a2 = q.W * q.Z
		R21 = a1 + a2
		R12 = a1 - a2
		a1 = q.X * q.Z
		a2 = q.W * q.Y
		R13 = a1 + a2
		R31 = a1 - a2
		a1 = q.Y * q.Z
		a2 = q.W * q.X
		R32 = a1 + a2
		R23 = a1 - a2
		R11 += R11
		R12 += R12
		R13 += R13
		R21 += R21
		R22 += R22
		R23 += R23
		R31 += R31
		R32 += R32
		R33 += R33

		// rotate m to earth frame and compute b
		a1 = R11*p.mx + R12*p.my + R13*p.mz
		a2 = R21*p.mx + R22*p.my + R23*p.mz
		bx = math.Sqrt(a1*a1 + a2*a2)
		bz = R31*p.mx + R32*p.my + R33*p.mz

		// compute J_g^T * f_g to compute qe
		f1 = R31 - p.ax
		f2 = R32 - p.ay
		f3 = R33 - p.az
		a1 = q.X * f3
		a2 = q.Y * f3
		qew = -q.Y*f1 + q.X*f2
		qex = q.Z*f1 + q.W*f2 - a1 - a1
		qey = -q.W*f1 + q.Z*f2 - a2 - a2
		qez = q.X*f1 + q.Y*f2

		// compute J_b^T * f_b to compute qe
		f1 = R11*bx + R31*bz - p.mx
		f2 = R12*bx + R32*bz - p.my
		f3 = R13*bx + R33*bz - p.mz
		a1 = q.W * bx
		a2 = q.X * bx
		a3 = q.Y * bx
		a4 = q.Z * bx
		a5 = q.W * bz
		a6 = q.X * bz
		a7 = q.Y * bz
		a8 = q.Z * bz
		qew += -a7*f1 + (a6-a4)*f2 + a3*f3
		qex += a8*f1 + (a3+a5)*f2 + (a4-a6-a6)*f3
		qey += (-a5-a3-a3)*f1 + (a2+a8)*f2 + (a1-a7-a7)*f3
		qez += (a6-a4-a4)*f1 + (a7-a1)*f2 + a2*f3

		// normalise qe
		norm = math.Sqrt(qew*qew + qex*qex + qey*qey + qez*qez)
		qew /= norm
		qex /= norm
		qey /= norm
		qez /= norm

		// compute omega_b

		// compute q_dot_omega
		a1 = -q.X*p.wx - q.Y*p.wy - q.Z*p.wz
		a2 = q.W*p.wx - q.Z*p.wy + q.Y*p.wz
		a3 = q.Z*p.wx + q.W*p.wy - q.X*p.wz
		a4 = -q.Y*p.wx + q.X*p.wy + q.W*p.wz
		a1 /= 2
		a2 /= 2
		a3 /= 2
		a4 /= 2

		// compute q_dot
		a1 -= p.beta * qew
		a2 -= p.beta * qex
		a3 -= p.beta * qey
		a4 -= p.beta * qez

		// compute q
		q.W += a1 * p.dt
		q.X += a2 * p.dt
		q.Y += a3 * p.dt
		q.Z += a4 * p.dt

		// normalise q
		norm = math.Sqrt(q.W*q.W + q.X*q.X + q.Y*q.Y + q.Z*q.Z)
		q.W /= norm
		q.X /= norm
		q.Y /= norm
		q.Z /= norm
	}
}
