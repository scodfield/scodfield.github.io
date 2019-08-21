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
