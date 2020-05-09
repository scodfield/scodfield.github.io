1. 只出现一次的数字
   异或解法,erlang版本:
   single_number(NumberList) ->
     lists:foldl(fun(N,R) -> N bxor R end, 0, NumberList). 
   借助额外数组,golang版本:
   func singleNumber(nums []int) int {
	  numberMap := make(map[int]int)

 	  for i := 0; i < len(nums); i++ {
 		 count, ok := numberMap[nums[i]]
 		 if ok {
 			numberMap[nums[i]] = count + 1
 		 } else {
 			numberMap[nums[i]] = 1
 		 }
 	  }

     var result int
 	  for i := range(numberMap) {
 		 if numberMap[i] == 1 {
 			 result = i
 		 }
 	 }
    
     return result
   }
2. 求众数
   常规解法,golang版本:
   func majorityElement(nums []int) int {
    numberMap := make(map[int]int)

    for i := 0; i < len(nums); i++ {
 		count, ok := numberMap[nums[i]]
 		if ok {
 			numberMap[nums[i]] = count + 1
 		} else {
 			numberMap[nums[i]] = 1
 		}
 	}

    result, count := 0 , 0
    for i := range(numberMap) {
    	if numberMap[i] >= len(nums)/2 && numberMap[i] >= count {
    		result, count = i, numberMap[i]
    	}
     }

     return result
  }
  摩尔投票法(Moore Voting),前提是必定存在出现次数过半的数,
  摩尔投票法先假定第一个元素过半数,计数器设为1,比较下一个数和次数是否相同,若相同则计数器+1,反之-1,之后判断计数器是否为0,若为0,则将下一个数设为候选
  众数,以此类推,直至遍历完整个数组,当前候选即为要求的众数
  摩尔投票参考:https://www.zhihu.com/question/49973163
3. 搜索二维矩阵ii,由于行列均满足有序,最开始的想法是每一行执行二分查找,或者将二维数据归并排序,合并成一维数组,再执行二分查找,时间复杂度不符合'高效'的
   要求,基于从左到右,从上到下升序的前提,可以从左下角开始比较,若当前值等于目标值,则返回true,若当前值小于目标值,则右移,反之则左移,直到矩阵的右上角
   以下为golang版本的实现:
   func searchMatrix(matrix [][]int, target int) bool {
	for row, col := len(matrix) - 1, 0; row >= 0 && col <= len(matrix[0]) -1; {
		if matrix[row][col] == target {
			return true
		} else if matrix[row][col] > target {
			row -= 1
		} else {
			col += 1
		}
	}

	return false
    }
4. 两数之和,常规思路是双层循环,找到目标之后直接返回,时间复杂度高,其次是借助使用hash结构的数据类型缓存中间计算结果,一下为使用map的golang版本:
	func twoSum(nums []int, target int) []int {
    var result []int
    var diffs map[int]int
    diffs = make(map[int]int)

    for i := 0; i < len(nums); i++ {
    	diff := target - nums[i]
    	if index, ok := diffs[diff]; ok {
    		if index > i {
    			result = append(result,i)
    			result = append(result,index)
    			return result
    		} else {
    			result = append(result,index)
    			result = append(result,i)
    			return result
    		}
    	}

    	diffs[nums[i]] = i
     }
    
      return result
   }
5. 两数相加,常规思路是先判断两个链表是否为nil,循环两个不为空的链表,最后在判断进位是否大于0,大于0则需要多一个节点,另一种思路是用递归来实现
   以下代码为golang版的递归实现:
   func addTwoNumbers2(l1 *ListNode, l2 *ListNode) *ListNode { 
    return addtwo(l1,l2,0)
   }

   func addtwo(l1 *ListNode, l2 *ListNode, carry int) *ListNode {
	x, y := 0, 0
	if l1 != nil {
		x = l1.Val
	}

	if l2 != nil {
		y = l2.Val
	}

	sum := x + y + carry

       var newNode ListNode
       newNode.Val = sum % 10
       carry = sum / 10
       if l1.Next == nil || l2.Next == nil {
    	   if l1.Next != nil {
    		var l2Node ListNode
    		l2Node.Val = 0
    		l2.Next = &l2Node
    		newNode.Next = addtwo(l1.Next,l2.Next,carry)
    	   } else if l2.Next != nil {
    		var l1Node ListNode
    		l1Node.Val = 0
    		l1.Next = &l1Node
    		newNode.Next = addtwo(l1.Next,l2.Next,carry)
    	   } else if carry > 0 {
    		var lastNode ListNode
    		lastNode.Val = carry
    		newNode.Next = &lastNode
    	   }
       } else {
    	  newNode.Next = addtwo(l1.Next,l2.Next,carry)
       }

      return &newNode
   }
6. 无重复字符的最长子串,常规思路还是双层循环,统计每个字符可以构造的最长无重复子串,最后返回最长子串的长度即可,贴一个自己的实现:
    func lengthOfLongestSubstring(s string) int {
    var result int = 0

    if s == "" {
    	return 0
    } else if len(s) == 1 {
    	return 1
    }

    var count int = 1
    for i := 0; i < len(s) - 1; i++ {
    	temp := make(map[uint8]int)
    	temp[s[i]] = 1
    	count = 1
    	for j := i+1; j < len(s); j++ {
    		_, ok := temp[s[j]]
    		if ok {
    			break
    		} else {
    			count ++
    			temp[s[j]] = 1
    		}
    	 }

    	 if result < count {
    		result = count
    	 }
      }

      return result
   }
     上述代码提交后,发现不关是时间还是空间都很垃圾,参考了评论,发现有更好的思路,代码如下:
     func lengthOfLongestSubstring(s string) int {
	i, j, max := 0, 0, 0

	for ; i < len(s); i++ {
		for x := j; x < i; x++ {
			if s[x] == s[i] {
				j = x + 1
				break
			}
		}
		fmt.Println("i: ",i,", j: ",j)
		if i-j+1 > max {
			max = i - j + 1
		}
	}
	return max
     }
     由代码可知解题思路为求当前下标到最大下标之间的不重复子串,一旦出现重复字符,则当前子串查找结束,判断是否是当前最长,更新max的值,同时更改当前下标
     为当前重复字符的下一个下标值,i,j,max初始值及相关计算方式可以兼顾到s的不同情况(为空或者只有一个字符等)
7. 寻找两个有序数组的中位数,题目要去log(m+n),目前只实现了m+n的时间复杂度,首先合并两个有序数组,再根据奇偶计算中位数,实现如下:
   func findMedianSortedArrays(nums1 []int, nums2 []int) float64 {
    // nums3 := merge(nums1,len(nums1),nums2,len(nums2))
    nums3 := merge_sorted(nums1,nums2)
    fmt.Println("median nums3",nums3)
    if len(nums3) % 2 == 0 {
    		return float64(nums3[len(nums3)/2] + nums3[len(nums3)/2-1]) / float64(2)
    	} else {
    		return float64(nums3[len(nums3)/2])
    	}
   }

   func merge_sorted(nums1 []int, nums2 []int) []int {
	var nums3 []int
	i, j := 0, 0

	for i < len(nums1) && j < len(nums2) {
		if nums1[i] < nums2[j] {
			nums3 = append(nums3,nums1[i])
			i++
		} else if nums1[i] > nums2[j] {
			nums3 = append(nums3,nums2[j])
			j++
		} else {
			nums3 = append(nums3,nums2[j])
			nums3 = append(nums3,nums2[j])
			i++
			j++
		}
	}

	for ; i < len(nums1); i++ {
		nums3 = append(nums3,nums1[i])
	}
	for ; j < len(nums2); j++ {
		nums3 = append(nums3,nums2[j])
	}

	return nums3
    }
8. 二叉树的中序遍历,常规实现为递归,但时间和空间复杂度均为O(n),中序遍历的一个非常经典的算法是莫里斯(Morris)遍历,它可以做到时间O(n),空间O(1)
   莫里斯算法的关键在于当前节点的左子树中的最右节点,将该左右节点的右空指针指向当前节点,用于遍历时候的回退,在中序遍历序列中,该最右节点刚好是当前
   节点的前驱节点
   Morris参考:https://blog.csdn.net/koalacoco/article/details/60464975
9. 删除链表的倒数第N个节点,实现了常规思路,两次遍历,第一次遍历(O(L)),求出链表长度L,求出要删除的第(L-N+1)个节点,第二次遍历(O(L-N)),找到第L-N个
   节点,修改 (L-N).Next = (L-N).Next.Next即可,所以该解法的时间复杂度为O(L),空间复杂度为常量O(1),题目提示能否只遍历一遍,我的想法是,遍历一遍
   统计两个信息,其一是链表的长度,其二将链表每一个节点的地址依次放入数组,则遍历之后,只需常量时间复杂度即可完成对链表的修改,不过空间复杂度为O(L)
   看了下官方题解,通过双指针来实现只遍历一次,且时间复杂度为O(L),空间复杂度为O(1)
   上述双指针单次遍历题解中,除了双指针这个点,还有一个点是引入了一个哑节点,哑节点的引入极大的简化了一些极端情况,比如只有链表一个节点或者删除链表
   头部节点,上述自己实现的常规思路解法,还需要特殊处理这些极端情况,代码很不整洁
10. 删除排序链表中的重复元素II,基于三指针&哑指针,实现删除链表中重复元素,引入哑指针后,可有效处理[1,1],[1,2,2]等特殊情况,代码如下:
    func deleteDuplicates(head *ListNode) *ListNode {
      var dummy,pre,first,second *ListNode
      dummy = &ListNode{0,nil}
      pre,first,second = dummy,head,head

      for second != nil {
    	  for (second.Next != nil) && (second.Next.Val == first.Val) {
    		second = second.Next
          }
    	  if second == first {
    	    pre.Next = first
    	    pre = first
    	  } 
    	  first,second = second.Next,second.Next
       }
       pre.Next = first
       
       return dummy.Next
    }
11. 分隔链表,构造一个以目标x为Val的节点,遍历原始链表,重新构造链表,结果符合题目要求,但是不给过,代码如下:
    func partition(head *ListNode, x int) *ListNode {
      var mid,front,tail,temp *ListNode
      mid = &ListNode{x,nil}
      front,tail = mid,mid

      for head != nil {
    	  temp = head.Next
    	  if head.Val < x {
    		head.Next = front
    		front = head
    	  } else if head.Val > x {
    		tail.Next = head
    		tail = head
    	  }
    	  head = temp
      }
      tail.Next = nil

      return front
    }
12. 最长回文子串,最开始的想法是这个问题和求字符串的最长不重复子串类似,所以考虑使用动态规划,如果一个子串str[i,j](以i,j开始&结尾)为回文,那么
    str[i+1,j-1]肯定也是回文,设定p[i,j]=0表示子串[i,j]不是回文子串,p[i,j]=1表示子串[i,j]为回文,则状态转移方程为: p[i,j] = p[i+1,j-1] 
    (str[i] == str[j]); p[i,j] = 0 (str[i] != str[j]), 实现代码如下:
    func longestPalindrome(s string) string {
	    var mem [20][20]int
	    left,right := 0,0

	    if len(s) == 0 {
		return ""
	    } else if len(s) == 1 {
		return s
	    }

	    for i := 0; i < len(s); i++ {
		mem[i][i] = 1
		if (i < len(s)-1) && (s[i] == s[i+1]) {
			mem[i][i+1] = 1
			left,right = i,i+1
		}
	    }

	    j := 0
	    for l := 3; l <= len(s); l++ {
		for i := 0; i+l-1 < len(s); i++ {
			j = i+l-1
			if (s[i]==s[j]) && (mem[i+1][j-1]==1) {
				mem[i][j] = 1
				left,right = i,j
			}
		}
	    }

	    return s[left:right+1]
	}
    动态规划算法求解最长回文子串,时间复杂度为O(N^2),同时需要借助二维数组保存中间计算结果,空间复杂度为O(N^2),从复杂度上看,该求解算法效率并不高
    看了下官方题解,有时间复杂度为O(N^2),空间复杂度为O(1)的中心扩展法,也有复杂度为O(N)的Manacher算法
    Manacher算法可参考:https://segmentfault.com/a/1190000003914228; https://cloud.tencent.com/developer/news/312855
13. N皇后
14. 最大子序和,典型的动态规划,代码如下:
    func maxSubArray(nums []int) int {
      if len(nums) <= 0 {
    	  return 0
      }
      if len(nums) <= 1 {
    	  return nums[0]
      }
      var curMax,preMax,i int
      curMax,preMax,i = nums[0],nums[0],1
      for ; i < len(nums); i++ {
    	  temp := preMax + nums[i]
    	  if temp > nums[i] {
    	  	preMax = temp
    	  } else {
    		preMax = nums[i]
    	  }
    	  if preMax > curMax {
    		curMax = preMax
    	  }
      }
      return curMax
   }
14. 买卖股票的最佳时机ii,扩展买卖股票的最佳时机i的解法,因为可以多次交易,所以记录当前这一手交易的当前最大获利&最低股价,
    若第i天出手交易获利小于当前已获得的最大获利,则当前这一手交易结束,从第i天开始新的一手交易,重置新的当前获利&最低股价,代码如下:
    func maxProfit(prices []int) int {
      if len(prices) <= 1 {
 		return 0   		
      }
      var profit,curMaxProfit,min,i int
      profit,curMaxProfit,min,i = 0,0,prices[0],1
      for ; i < len(prices); i++ {
 	temp := prices[i] - min
 	if temp >= curMaxProfit {
 		curMaxProfit = temp
 	} else {
 		profit = profit + curMaxProfit
 		curMaxProfit = 0
 		min = prices[i]
 	}
      }
      return profit + curMaxProfit
    }
14. 颠倒二进制位,与"只出现一次的数"中的位运算解法类似,本题也可以用位运算解决,代码如下:
    func reverseBits(num uint32) uint32 {
      var result uint32
      for i := 0; i < 32; i++ {
    	 result <<= 1
    	 result |= (num & 1)
    	 num >>= 1
      }
      return result
    }
15. 盛最多水的容器,最开始的思路是内外双层循环,计算以每一个坐标为最左边界时,与其右边各坐标所构成的容器容量,时间复杂度O(N^2),
    显然这种解法too native,看了官方题解可以采用"双指针法",从数组两端开始往中间移动,以左右指针所指向坐标值为高度,两坐标差为长度,
    双指针法每次移动较高度低一侧的指针,直到两指针相遇,虽然在移动过程中,长度在逐渐变小,但是高度有可能会变大,这可以保证在移动过程中
    新容器的容量变大,同样的之所以移动高度较低的一侧,也是因为长度在逐渐变小,如果移动高度较高一侧的指针,则高度取决于较低一侧,而长度
    在减少,那么新容器容量绝不会超过之前的容量,双指针法时间复杂度为O(N),代码如下:
    func maxArea(height []int) int {
      if len(height) <= 1 {
            return 0
	}
      var maxArea,left,right,tempArea int
      left, right = 0, len(height)-1
      for left < right {
    	if height[left] < height[right] {
    		tempArea = height[left] * (right - left)
    		left++
    	} else {
    	    tempArea = height[right] * (right - left)
    		right--
    	}
    	if maxArea < tempArea {
    	    maxArea = tempArea
    	 }
        }
        return maxArea
      }
16. 接雨水,最初的思路是将每一个坐标作为可盛放雨水的容器的bottom,计算基于该bottom的容器的容量,所计算的区域不止包括该坐标的容量,
    还可能会包含左右多个坐标(求取的区域是可能跨多个坐标的一个长方形),算法时间复杂度为O(N^2),需注意的是,在横向扩展坐标时,有可能
    遇到和当前计算的坐标相同高度的其它坐标,为避免重复计算,需借助额外的数组标记坐标是否已被同高度的其它坐标计算过,代码较长,时间
    复杂度较高,看了官方题解,双指针法实现O(N^2)复杂度,算法的思路和盛水容器类似,移动高度较低一端的指针,同时维护leftMax&rightMax,
    不过要注意的是,和第一种思路相比,计算的是当前坐标上所能存储的水的多少(并不是横跨多个坐标的水量),代码如下:
    func trap(height []int) int {
        if len(height) <= 2 {
    	    return 0
        }
        var total,left,right,leftMax,rightMax int
        left, right = 0, len(height) - 1
        leftMax, rightMax = height[left], height[right]
        for left < right {
    	    if height[left] < height[right] {
    		if height[left] >= leftMax {
    			leftMax = height[left]
    		} else {
    			total += leftMax - height[left]
    		}
    		left++
    	    } else {
    		if height[right] >= rightMax {
    			rightMax = height[right]
    		} else {
    			total += rightMax - height[right]
    		}
    		right--
    	    }
        }
        return total
    }
17.字符串转换整数,本题需要注意异常及边界条件(空字符串,只包含空格的字符串,+/-后续空格,32bite的最大/小值),代码如下:
	func myAtoi(str string) int {
	    maxInt :=  1 << 31 - 1
	    minInt := -1 << 31
	    strLen := len(str)
	    var result,i,sign int

	    // trim left
	    for (i < strLen) && (str[i] == ' ') {
		i++
	    } 

	    if (strLen <= 0) || (i >= strLen) || (!isDigit(str[i]) && !isSign(str[i])) {
		return 0
	    }

	    // signed
	    if str[i] == '+' {
		sign = 1
		i++
	    } else if str[i] == '-' {
		sign = -1
		i++
	    } else {
		sign = 1
	    }

	    // loop
	    for i < strLen && isDigit(str[i]) {
		result = result * 10 + int(str[i] - '0')
		if result > maxInt {
		    break
		}
		i++
	    }
	    result = result * sign

	    if result > maxInt {
		result = maxInt
	    } else if result < minInt {
		result = minInt
	    }

	    return result
	}

	func isDigit(c uint8) bool {
	    return '0' <= c && c <= '9'
	}

	func isSign(c uint8) bool {
	    return '+' == c || c == '-'
	}
