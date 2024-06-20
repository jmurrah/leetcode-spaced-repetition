package main

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/supabase-community/supabase-go"
)

type LeetCodeProblem struct {
	TitleSlug      string
	Difficulty     string
	Tags           []string
	CompletedDates []time.Time
	RepeatDate     time.Time
}

type EmailListResponse struct {
	EmailList []string `json:"email_list"`
}

func create_supabase_client() *supabase.Client {
	client, err := supabase.NewClient(os.Getenv("SUPABASE_URL"), os.Getenv("SUPABASE_KEY"), &supabase.ClientOptions{})
	if err != nil {
		fmt.Println("cannot initalize client", err)
	}
	return client
}

func add_username_and_email_to_database(leetcode_username string, new_email string) {
	client := create_supabase_client()
	table := os.Getenv("SUPABASE_TABLE")

	found_username, _, _ := client.From(table).Select("username", "", false).Eq("username", leetcode_username).Execute()

	if string(found_username) == "[]" {
		data := map[string]interface{}{
			"username":   leetcode_username,
			"email_list": []string{new_email},
			"problems":   []LeetCodeProblem{},
		}
		client.From(table).Insert(data, false, "Failure", "Success", "1").Execute()
		return
	}

	var old_emails []EmailListResponse
	var updatedEmails []string

	response, _, _ := client.From(table).Select("email_list", "", false).Eq("username", leetcode_username).Execute()
	json.Unmarshal(response, &old_emails)

	for _, old_email := range old_emails {
		updatedEmails = append(updatedEmails, old_email.EmailList...)
	}
	updatedEmails = append(updatedEmails, new_email)

	data := map[string]interface{}{
		"email_list": updatedEmails,
	}
	client.From(table).Update(data, "", "").Eq("username", leetcode_username).Execute()
}
