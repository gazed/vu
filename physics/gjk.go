// Copyright © 2024 Galvanized Logic Inc.

package physics

import (
	"log/slog"

	"github.com/gazed/vu/math/lin"
)

// gjk_Simplex;
type gjk_Simplex struct {
	a, b, c, d lin.V3
	num        uint32
}

// add_to_simplex
func add_to_simplex(simplex *gjk_Simplex, point lin.V3) {
	switch simplex.num {
	case 1:
		simplex.b = simplex.a
		simplex.a = point
	case 2:
		simplex.c = simplex.b
		simplex.b = simplex.a
		simplex.a = point
	case 3:
		simplex.d = simplex.c
		simplex.c = simplex.b
		simplex.b = simplex.a
		simplex.a = point
	default:
		slog.Error("add_to_simplex")
	}
	simplex.num += 1
}

// triple_cross
func triple_cross(a, b, c lin.V3) (tc lin.V3) {
	tc.Cross(&a, &b)
	tc.Cross(&tc, &c)
	return tc
}

// do_simplex_2
func do_simplex_2(simplex *gjk_Simplex, direction *lin.V3) bool {
	a := simplex.a // the last point added
	b := simplex.b
	ao := lin.NewV3().Neg(&a)
	ab := lin.NewV3().Sub(&b, &a)
	if ab.Dot(ao) >= 0.0 {
		simplex.a = a
		simplex.b = b
		simplex.num = 2
		*direction = triple_cross(*ab, *ao, *ab)
	} else {
		simplex.a = a
		simplex.num = 1
		*direction = *ao
	}
	return false
}

// do_simplex_3
func do_simplex_3(simplex *gjk_Simplex, direction *lin.V3) bool {
	a := simplex.a // the last point added
	b := simplex.b
	c := simplex.c
	ao := lin.NewV3().Neg(&a)
	ab := lin.NewV3().Sub(&b, &a)
	ac := lin.NewV3().Sub(&c, &a)
	abc := lin.NewV3().Cross(ab, ac)

	if lin.NewV3().Cross(abc, ac).Dot(ao) >= 0.0 {
		if ac.Dot(ao) >= 0.0 {
			// AC region
			simplex.a = a
			simplex.b = c
			simplex.num = 2
			*direction = triple_cross(*ac, *ao, *ac)
		} else {
			if ab.Dot(ao) >= 0.0 {
				// AB region
				simplex.a = a
				simplex.b = b
				simplex.num = 2
				*direction = triple_cross(*ab, *ao, *ab)
			} else {
				// A region
				simplex.a = a
				*direction = *ao
			}
		}
	} else {
		if lin.NewV3().Cross(ab, abc).Dot(ao) >= 0.0 {
			if ab.Dot(ao) >= 0.0 {
				// AB region
				simplex.a = a
				simplex.b = b
				simplex.num = 2
				*direction = triple_cross(*ab, *ao, *ab)
			} else {
				// A region
				simplex.a = a
				*direction = *ao
			}
		} else {
			if abc.Dot(ao) >= 0.0 {
				// ABC region ("up")
				simplex.a = a
				simplex.b = b
				simplex.c = c
				simplex.num = 3
				*direction = *abc
			} else {
				// ABC region ("down")
				simplex.a = a
				simplex.b = c
				simplex.c = b
				simplex.num = 3
				*direction = *(abc.Neg(abc))
			}
		}
	}
	return false
}

// do_simplex_4
func do_simplex_4(simplex *gjk_Simplex, direction *lin.V3) bool {
	a := simplex.a // the last point added
	b := simplex.b
	c := simplex.c
	d := simplex.d

	ao := lin.NewV3().Neg(&a)
	ab := lin.NewV3().Sub(&b, &a)
	ac := lin.NewV3().Sub(&c, &a)
	ad := lin.NewV3().Sub(&d, &a)
	abc := lin.NewV3().Cross(ab, ac)
	acd := lin.NewV3().Cross(ac, ad)
	adb := lin.NewV3().Cross(ad, ab)

	plane_information := uint8(0)
	if abc.Dot(ao) >= 0.0 {
		plane_information |= 0x1
	}
	if acd.Dot(ao) >= 0.0 {
		plane_information |= 0x2
	}
	if adb.Dot(ao) >= 0.0 {
		plane_information |= 0x4
	}
	switch plane_information {
	case 0x0:
		// Intersection
		return true
	case 0x1:
		// Triangle ABC
		if lin.NewV3().Cross(abc, ac).Dot(ao) >= 0.0 {
			if ac.Dot(ao) >= 0.0 {
				// AC region
				simplex.a = a
				simplex.b = c
				simplex.num = 2
				*direction = triple_cross(*ac, *ao, *ac)
			} else {
				if ab.Dot(ao) >= 0.0 {
					// AB region
					simplex.a = a
					simplex.b = b
					simplex.num = 2
					*direction = triple_cross(*ab, *ao, *ab)
				} else {
					// A region
					simplex.a = a
					*direction = *ao
				}
			}
		} else {
			if lin.NewV3().Cross(ab, abc).Dot(ao) >= 0.0 {
				if ab.Dot(ao) >= 0.0 {
					// AB region
					simplex.a = a
					simplex.b = b
					simplex.num = 2
					*direction = triple_cross(*ab, *ao, *ab)
				} else {
					// A region
					simplex.a = a
					*direction = *ao
				}
			} else {
				// ABC region
				simplex.a = a
				simplex.b = b
				simplex.c = c
				simplex.num = 3
				*direction = *abc
			}
		}
	case 0x2:
		// Triangle ACD
		if lin.NewV3().Cross(acd, ad).Dot(ao) >= 0.0 {
			if ad.Dot(ao) >= 0.0 {
				// AD region
				simplex.a = a
				simplex.b = d
				simplex.num = 2
				*direction = triple_cross(*ad, *ao, *ad)
			} else {
				if ac.Dot(ao) >= 0.0 {
					// AC region
					simplex.a = a
					simplex.b = c
					simplex.num = 2
					*direction = triple_cross(*ab, *ao, *ab)
				} else {
					// A region
					simplex.a = a
					*direction = *ao
				}
			}
		} else {
			if lin.NewV3().Cross(ac, acd).Dot(ao) >= 0.0 {
				if ac.Dot(ao) >= 0.0 {
					// AC region
					simplex.a = a
					simplex.b = c
					simplex.num = 2
					*direction = triple_cross(*ac, *ao, *ac)
				} else {
					// A region
					simplex.a = a
					*direction = *ao
				}
			} else {
				// ACD region
				simplex.a = a
				simplex.b = c
				simplex.c = d
				simplex.num = 3
				*direction = *acd
			}
		}
	case 0x3:
		// Line AC
		if ac.Dot(ao) >= 0.0 {
			simplex.a = a
			simplex.b = c
			simplex.num = 2
			*direction = triple_cross(*ac, *ao, *ac)
		} else {
			simplex.a = a
			simplex.num = 1
			*direction = *ao
		}
	case 0x4:
		// Triangle ADB
		if lin.NewV3().Cross(adb, ab).Dot(ao) >= 0.0 {
			if ab.Dot(ao) >= 0.0 {
				// AB region
				simplex.a = a
				simplex.b = b
				simplex.num = 2
				*direction = triple_cross(*ab, *ao, *ab)
			} else {
				if ad.Dot(ao) >= 0.0 {
					// AD region
					simplex.a = a
					simplex.b = d
					simplex.num = 2
					*direction = triple_cross(*ad, *ao, *ad)
				} else {
					// A region
					simplex.a = a
					*direction = *ao
				}
			}
		} else {
			if lin.NewV3().Cross(ad, adb).Dot(ao) >= 0.0 {
				if ad.Dot(ao) >= 0.0 {
					// AD region
					simplex.a = a
					simplex.b = d
					simplex.num = 2
					*direction = triple_cross(*ad, *ao, *ad)
				} else {
					// A region
					simplex.a = a
					*direction = *ao
				}
			} else {
				// ADB region
				simplex.a = a
				simplex.b = d
				simplex.c = b
				simplex.num = 3
				*direction = *adb
			}
		}
	case 0x5:
		// Line AB
		if ab.Dot(ao) >= 0.0 {
			simplex.a = a
			simplex.b = b
			simplex.num = 2
			*direction = triple_cross(*ab, *ao, *ab)
		} else {
			simplex.a = a
			simplex.num = 1
			*direction = *ao
		}
	case 0x6:
		// Line AD
		if ad.Dot(ao) >= 0.0 {
			simplex.a = a
			simplex.b = d
			simplex.num = 2
			*direction = triple_cross(*ad, *ao, *ad)
		} else {
			simplex.a = a
			simplex.num = 1
			*direction = *ao
		}
	case 0x7:
		// Point A
		simplex.a = a
		simplex.num = 1
		*direction = *ao
	}
	return false
}

// do_simplex
func do_simplex(simplex *gjk_Simplex, direction *lin.V3) bool {
	switch simplex.num {
	case 2:
		return do_simplex_2(simplex, direction)
	case 3:
		return do_simplex_3(simplex, direction)
	case 4:
		return do_simplex_4(simplex, direction)
	}
	return false
}

// gjk_collides
func gjk_collides(collider1, collider2 *collider, _simplex *gjk_Simplex) bool {
	var simplex gjk_Simplex
	simplex.a = support_point_of_minkowski_difference(collider1, collider2, lin.V3{0, 0, 1})
	simplex.num = 1
	direction := lin.NewV3().Scale(&simplex.a, -1.0)
	for i := 0; i < 100; i++ {
		next_point := support_point_of_minkowski_difference(collider1, collider2, *direction)
		if next_point.Dot(direction) < 0.0 {
			// No intersection.
			return false
		}
		add_to_simplex(&simplex, next_point)
		if do_simplex(&simplex, direction) {
			// Intersection.
			if _simplex != nil {
				*_simplex = simplex
			}
			return true
		}
	}
	// slog.Warn("GJK did not converge")
	return false
}
