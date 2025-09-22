package main

import (
	"encoding/json"
	"log"
	"net/http"
)

// 메시지 구조체
type Message struct {
	// Text라는 이름의 문자열 변수를 하나 가진다.
	// 태그: json으로 바꿀 때, Text라는 변수명을 message라는 이름으로 바꿔주세요
	Text string `json:"message"`
}

func main() {
	// /button1 주소에 누가 접속했을 때, func 안쪽 코드를 실행함
	// 받는 요청 정보: r / 보낼 응답 정보: w
	http.HandleFunc("/button1", func(w http.ResponseWriter, r *http.Request) {
		msg := Message{Text: "Button 1 Clicked"}
		// 클라이언트한테 이거 json이라고 알려줌
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(msg)
	})

	http.HandleFunc("/button2", func(w http.ResponseWriter, r *http.Request) {
		msg := Message{Text: "Button 2 Clicked"}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(msg)
	})

	log.Println("Server Started!")

	// ListenAndServe: 8080 포트로 들어오는 모든 요청을 기다리며 처리함
	// 그러면서 치명적인 문제가 생기면 터미널에 출력하고 서버 종료
	log.Fatal(http.ListenAndServe(":8080", nil))
}