package main

import (
	"fmt"
	"math"
	"strconv"
)

var charMap = map[rune]rune{')': '(', '}': '{', ']': '['}

func FibonacciIterative(n int) int {
	// Функція вираховує і повертає n-не число Фібоначчі
	// Імплементація без використання рекурсії
	if n < 2 {
		return n
	}
	a := 0
	b := 1
	for i := 2; i <= n; i++ {
		a, b = b, a+b
	}
	return b
}

func FibonacciRecursive(n int) int {
	// Функція вираховує і повертає n-не число Фібоначчі
	// Імплементація з використанням рекурсії
	if n < 2 {
		return n
	} else {
		return FibonacciRecursive(n-1) + FibonacciRecursive(n-2)
	}
}

func IsPrime(n int) bool {
	// Функція повертає `true` якщо число `n` - просте.
	// Інакше функція повертає `false`
	if n <= 1 {
		return false
	} else if n == 2 {
		return true
	} else if n%2 == 0 {
		return false
	}
	sqrtN := int(math.Sqrt(float64(n)))
	for i := 3; i <= sqrtN; i++ {
		if n%i == 0 {
			return false
		}
	}
	return true
}

func IsBinaryPalindrome(n int) bool {
	// Функція повертає `true` якщо число `n` у бінарному вигляді є паліндромом
	// інакша функція повертає `false`
	//
	// Приклади:
	// Число 7 (111) - паліндром, повертаємо `true`
	// Число 5 (101) - паліндром, повертаємо `true`
	// Число 6 (110) - не є паліндромом, повертаємо `false`
	binaryRepresentation := strconv.FormatInt(int64(n), 2)
	length := len(binaryRepresentation)
	for i := 0; i < length/2; i++ {
		if binaryRepresentation[i] != binaryRepresentation[length-i-1] {
			return false
		}
	}
	return true
}

func ValidParentheses(s string) bool {
	// Функція повертає `true` якщо у вхідній стрічці дотримані усі правила використання дужок
	// Правила:
	// 1. Допустимі дужки `(`, `[`, `{`, `)`, `]`, `}`
	// 2. У кожної відкритої дужки є відповідна закриваюча дужка того ж типу
	// 3. Закриваючі дужки стоять у правильному порядку
	//    "[{}]" - правильно
	//    "[{]}" - не правильно
	// 4. Кожна закриваюча дужка має відповідну що відкриває дужку
	stack := make([]rune, 0)
	for _, char := range s {
		if char == '(' || char == '{' || char == '[' {
			stack = append(stack, char)
		} else if char == ')' || char == '}' || char == ']' {
			if charMap[char] != stack[len(stack)-1] {
				return false
			} else {
				stack = stack[:len(stack)-1]
			}
		}
	}
	if len(stack) > 0 {
		return false
	}
	return true
}

func Increment(num string) int {
	// Функція на вхід отримує стрічку яка складається лише з символів `0` та `1`
	// Тобто стрічка містить певне число у бінарному вигляді
	// Потрібно повернути число на один більше
	inputInt, err := strconv.ParseInt(num, 2, 64)
	if err != nil {
		fmt.Printf("Invalid number %v\n", num)
		return 0
	}
	outputInt := inputInt + 1
	return int(outputInt)
}

func main() {
	fmt.Println(FibonacciIterative(6))                                                // 8
	fmt.Println(FibonacciRecursive(6))                                                // 8
	fmt.Println(IsPrime(11))                                                          // true
	fmt.Println(IsBinaryPalindrome(5))                                                // true
	fmt.Println(ValidParentheses("func() { return fmt.Println(len([]int{1,2,3}))}"))  // true
	fmt.Println(ValidParentheses("func() { return fmt.Println(len)([]int{1,2,3}))}")) // false
	fmt.Println(Increment("111"))                                                     // 8
	fmt.Println(Increment("101"))                                                     // 6
	fmt.Println(Increment("110"))                                                     // 7
}
