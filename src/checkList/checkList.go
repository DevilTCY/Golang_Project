package checkList

import (
	"math"
)

/*定义最小数列长度*/
const MINLINGTH = 3

/*利用等差数列性质 2 * An = A(n-1) + A(n+1) 判断数列是否为等差数列，并计算公差*/
func isArithmeticSequences(list []int, length int, d *int) bool {
	if length < MINLINGTH {
		return false
	}

	for i := 1; i < length-1; i++ {
		if list[i]*2 != list[i-1]+list[i+1] {
			return false
		}
	}
	*d = list[1] - list[0]
	return true
}

/*利用等比数列性质 An^2 = A(n+1) * A(n-1) 判断数列是否为等比数列*/
func isGeometricSequences(list []int, length int, flage, q *int) bool {
	if length < MINLINGTH {
		return false
	}

	for i := 1; i < length-1; i++ {
		if list[i]*list[i] != list[i+1]*list[i-1] {
			return false
		}
	}

	/*公比不能为0*/
	tmp := float64(list[1]) / float64(list[0])
	if tmp == 0 {
		return false
	}

	/*判断公比是否为小数*/
	_, frac := math.Modf(tmp)
	if 0 != frac {
		*flage = 1
		return false
	}
	*q = int(tmp)
	return true
}

/*求x的平方根：使用内置函数math.Sqrt获取float64平方根，判断是否有小数*/
func mySqrt(x int) int {
	ans := math.Sqrt(float64(x))
	/*使用math.Modf函数获取ans整数、小数部分*/
	_, frac_num := math.Modf(ans)
	if 0 != frac_num {
		return -1
	}
	return int(ans)
}

/*开方等差数列*/
func isSqrtSequences(list []int, length int, next_val *int) bool {
	if length < MINLINGTH {
		return false
	}

	/*先将list中每个元素开方,取正直*/
	tmplist := []int{}
	for _, val := range list {
		tmp := mySqrt(val)
		if tmp == -1 {
			return false
		}
		tmplist = append(tmplist, tmp)
	}

	/*在判断开方后的tmplist是否为等差数列*/
	d := 0
	if !isArithmeticSequences(tmplist, len(tmplist), &d) {
		return false
	}

	/*计算下一项值*/
	*next_val = myPowerf(tmplist[len(tmplist)-1]+d, 2)
	return true
}

/*求x的y次幂*/
func myPowerf(x, y int) int {
	ans := 1
	for y != 0 {
		if y%2 == 1 {
			ans *= x
		}
		x *= x
		y /= 2
	}
	return ans
}

/*组合嵌套*/
/*对数列每次操作记录*/
const (
	DEM      int = 1 //求差
	QUOTIENT int = 2 //求商
	SQRT     int = 3 //开方
)

/*判断数列能否进行求差、求商、开方
proce_type：数列操作类型：1 -- 求差；2 -- 求商；3 -- 开方*/
func listProcess(list []int, length, proce_type int, resList *[]int) bool {
	tmpList := []int{}
	if proce_type == DEM { /*数列求差*/
		for i := 1; i < length; i++ {
			tmpList = append(tmpList, list[i]-list[i-1])
		}
		if len(tmpList) < MINLINGTH {
			return false
		}

		*resList = tmpList
		return true
	} else if proce_type == QUOTIENT { /*数列求商*/
		var tmp float64
		for i := 1; i < length; i++ {
			/*判断除数是否为0*/
			if list[i-1] == 0 {
				return false
			}

			tmp = float64(list[i]) / float64(list[i-1])
			if tmp == 0 {
				return false
			}
			/*判断商是否有小数*/
			_, frac := math.Modf(tmp)
			if 0 != frac {
				return false
			}

			tmpList = append(tmpList, int(tmp))
		}

		if len(tmpList) < MINLINGTH {
			return false
		}

		*resList = tmpList
		return true
	} else if proce_type == SQRT {
		for _, val := range list {
			tmp := mySqrt(val)
			if tmp == -1 {
				return false
			}
			tmpList = append(tmpList, tmp)
		}

		if len(tmpList) < MINLINGTH {
			return false
		}

		*resList = tmpList
		return true
	}

	return false
}

/*ifx：记录已做过判断的数列，防止程序死循环*/
var lastList [][]int

func myEqual(x, y []int) bool {
	if len(x) != len(y) {
		return false
	}

	if (x == nil) != (y == nil) {
		return false
	}

	for k, val := range x {
		if val != y[k] {
			return false
		}
	}
	return true
}

func isHandled(list []int) bool {
	for _, val := range lastList {
		if myEqual(list, val) {
			return true
		}
	}
	return false
}
func isMultSequences(list []int, val *int, ok *bool) {
	if len(list) < MINLINGTH {
		return
	}

	/*fix:当前数列是否已判断：序列已经判断，说明触发了循环条件，直接退出，否则会造成死循环*/
	if isHandled(list) {
		return
	}
	lastList = append(lastList, list)

	tmp := 0
	tmpList := list

	/*判断差值是否构成等差数列*/
	if isArithmeticSequences(tmpList, len(tmpList), &tmp) {
		*val = tmpList[len(tmpList)-1] + tmp
		*ok = true
		return
	}

	/*判断是否构成等比数列*/
	flage := 0
	if isGeometricSequences(tmpList, len(tmpList), &flage, &tmp) {
		*val = tmpList[len(tmpList)-1] * tmp
		*ok = true
		return
	}

	/*判断是否构成开方等差数列*/
	if isSqrtSequences(tmpList, len(tmpList), &tmp) {
		*val = tmp
		*ok = true
		return
	}

	/*当前差值数列不构成任何有规律的数列，递归*/
	/*原始数列求差*/
	tmpList = tmpList[0:0]
	if listProcess(list, len(list), DEM, &tmpList) {
		isMultSequences(tmpList, val, ok)
		if *ok == true {
			*val = list[len(list)-1] + (*val)
			return
		} else {
			tmpList = tmpList[0:0]
			/*求商*/
			if listProcess(list, len(list), QUOTIENT, &tmpList) {
				isMultSequences(tmpList, val, ok)
				if *ok == true {
					*val = list[len(list)-1] * (*val)
					return
				}
			} else {
				tmpList = tmpList[0:0]
				/*开方*/
				if listProcess(list, len(list), SQRT, &tmpList) {
					isMultSequences(tmpList, val, ok)
					if *ok == true {
						*val = myPowerf(*val, 2)
						return
					}
				} else {
					return
				}
			}
		}
	}
}

/*
检查数组规律：
当数列同时满足多种规律时，同一级之间按照 等差数列 > 等比数列 > 开方等差数列 的优先级来计算
*/
func checkList(list []int) (val int, ok bool) {
	length := len(list)
	if length < MINLINGTH {
		return 0, false
	}

	/*等差数列：递增和递减*/
	d := 0
	if isArithmeticSequences(list, length, &d) {
		return list[length-1] + d, true
	}

	/*等比数列：公比可以为负整数，但不能为小数*/
	q := 0
	flage := 0
	if isGeometricSequences(list, length, &flage, &q) {
		return list[length-1] * q, true
	}
	/*判断公比是否为小数*/
	if flage == 1 {
		return 0, false
	}

	/*开方等差数列： 数组每一项可开方，且平方根是整数，平方根组成等差数列*/
	next_val := 0
	if isSqrtSequences(list, length, &next_val) {
		return next_val, true
	}

	/*数列任意嵌套*/
	ok = false
	isMultSequences(list, &val, &ok)
	lastList = lastList[0:0]
	return
}
