#!/usr/bin/env python3
import z3, re

def parse(a: str) -> list:
    vals, vels = a.split(" @ ")
    vals = [int(i) for i in vals.split(", ")]
    vels = [int(i) for i in vels.split(", ")]
    return vals+vels


with open("./input1.txt", "r") as f:
    h = f.readlines()
    h = [i.strip() for i in h]
    h = [parse(i) for i in h]

xi, yi, zi, vxi, vyi, vzi = z3.Ints( "xi yi zi vxi vyi vzi" )
times = [ z3.Int("t"+str(i)) for i in range(len(h)) ]

solver = z3.Solver()  # This feels wrong :(
for i, (x, y, z, vx, vy, vz) in enumerate(h):
    solver.add( x + vx * times[i] == xi + vxi * times[i] )
    solver.add( y + vy * times[i] == yi + vyi * times[i] )
    solver.add( z + vz * times[i] == zi + vzi * times[i] )
solver.check()

print(f"Part 2: {solver.model().evaluate(xi + yi + zi)}")
