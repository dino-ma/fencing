package fencing

import "math"

//围栏
type Fencing struct {
	Polygon [][2]float64 `json:"polygon"` //多边形的一组坐标点
}

//FenceNew 初始化对象
func FenceNew(polygon [][2]float64) *Fencing {
	return &Fencing{Polygon: polygon}
}

//rayCrossesSegment 向两个相交的节点发射射线，看是否会相交
func (f *Fencing) rayPassIntersect(point [2]float64, a [2]float64, b [2]float64) bool {
	//获取点的经纬度
	px, py := point[0], point[1]
	//获取两个坐标点的经纬度
	ax, ay := a[0], a[1]
	bx, by := b[0], b[1]
	//如果AY大于BY则做转换
	if ay > by {
		ax, ay = b[0], b[1]
		bx, by = a[0], a[1]
	}
	// alter longitude to cater for 180 degree crossings
	if px < 0 {
		px += 360
	}
	if ax < 0 {
		ax += 360
	}
	if bx < 0 {
		bx += 360
	}
	//如果要验证的点与A B两个坐标的Y轴重合 则给其增加一定高度
	if py == ay || py == by {
		py += 0.00000001
	}
	//点P的Y 不再A B点Y的范围上 或者 点X 大于 AB点最大的X 则无相交
	if (py > by || py < ay) || (px > math.Max(ax, bx)) {
		return false
	}
	//因为是向右发射射线，如果点P的X小于 AB点中最小的X则一定会相交
	if px < math.Min(ax, bx) {
		return true
	}

	var red, blue float64
	if ax != bx {
		red = (by - ay) / (bx - ax)
	} else {
		red = math.Inf(0)
	}

	if ax != px {
		blue = (py - ay) / (px - ax)
	} else {

		blue = math.Inf(0)
	}

	return blue >= red

}

//IsContainer 判断走标点是否在围栏内
func (f *Fencing) IsContainer(point [2]float64) bool {
	crossings := 0
	polygon := f.Polygon
	//最低也需要三个点
	if len(polygon) < 3 {
		return false
	}
	for i := 0; i < len(polygon); i++ {
		currentPoint := polygon[i]
		j := i + 1
		if j >= len(polygon) {
			j = 0
		}
		nextPoint := polygon[j]
		//根据我的坐标，向临近的两个坐标点发射射线，判断是否会相交
		if f.rayPassIntersect(point, currentPoint, nextPoint) {
			//返回BOOL则代表有一条射线相交
			crossings += 1
		}
	}

	//如果为奇数，Q在多边形内；如果为偶数，Q在多边形外
	return crossings%2 == 1
}
