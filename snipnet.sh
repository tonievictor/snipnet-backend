#!/usr/bin/env bash

signup() {
	echo "Please enter an email address:"
	read email

	echo "Please enter a username"
	read username

	echo "Please enter your password"
	read -s password

	res=$(curl -s -X POST -d "{\"username\": \"$username\", \"password\": \"$password\", \"email\": \"$email\"}" localhost:3000/signup)
	status=$(echo "$res" | jq -r '.status')

	if [[ "$status" = true ]]; then
		echo "Account created successfully"
		return
	else
		message=$(echo "$res" | jq -r '.message')
		echo "$message"
		return
	fi
}

signin() {
	echo "Please enter a username"
	read username

	echo "Please enter your password"
	read -s password

	res=$(curl -s -X POST -d "{\"username\": \"$username\", \"password\": \"$password\"}" localhost:3000/signin)
	status=$(echo "$res" | jq -r '.status')

	if [[ "$status" = true ]]; then
		auth_token=$(echo "$res" | jq -r '.data.auth_token')
		echo "$auth_token" > ~/.snipnetauth
		echo "Signed in successfully"
		return
	else
		message=$(echo "$res" | jq -r '.message')
		echo "$message"
		return
	fi
}

getsnippet() {
	echo "Please enter the snippet id"
	read id

	res=$(curl -s localhost:3000/snippets/$id)
	status=$(echo "$res" | jq -r '.status')

	if [[ "$status" = true ]]; then
		data=$(echo "$res" | jq -r '.data' | jq)
		echo "$data"
		return
	else
		message=$(echo "$res" | jq -r '.message')
		echo "$message"
		return
	fi
}

getsnippets() {
	auth_token=$(cat ~/.snipnetauth 2>> /dev/null)
	if [[ "$(echo "$?")" -ne 0 ]]; then
		echo "You are not logged in, please do that and try again."
		return
	fi
	res=$(curl -s -H "Authorization: Bearer $auth_token" localhost:3000/users/snippets)
	status=$(echo "$res" | jq -r '.status')

	if [[ "$status" = true ]]; then
		data=$(echo "$res" | jq -r '.data' | jq)
		echo "$data"
		return
	else
		message=$(echo "$res" | jq -r '.message')
		echo "$message"
		return
	fi
}

createsnippet() {
	auth_token=$(cat ~/.snipnetauth 2>> /dev/null)
	if [[ "$(echo "$?")" -ne 0 ]]; then
		echo "You are not logged in, please do that and try again."
		return
	fi
	echo "Please enter a title"
	read title

	echo "Please enter a description"
	read description

	echo "What language is this snippet?"
	read language

	echo "Provide the snippet here."
	read code

	echo "$code"

	res=$(curl -s -H "Authorization: Bearer $auth_token" -X POST -d "{\"title\": \"$title\", \"description\": \"$description\", \"language\": \"$language\", \"code\": \"$code\"}" localhost:3000/snippets)
	status=$(echo "$res" | jq -r '.status')

	if [[ "$status" = true ]]; then
		data=$(echo "$res" | jq -r '.data' | jq)
		echo "Snippet created successfully"
		echo "$data"
		return
	else
		message=$(echo "$res" | jq -r '.message')
		echo "$message"
		return
	fi
}

deletesnippet() {
	auth_token=$(cat ~/.snipnetauth 2>> /dev/null)
	if [[ "$(echo "$?")" -ne 0 ]]; then
		echo "You are not logged in, please do that and try again."
		return
	fi
	echo "Please enter the snippet id"
	read id

	res=$(curl -s -I -H "Authorization: Bearer $auth_token" -X DELETE localhost:3000/snippets/$id)
	status_code=$(echo "$res" | head -n 1 | awk '{print $2}')

	if [[ "$status_code" -eq 204 ]]; then
		echo "Snippet deleted successfully"
		return
	elif [[ "$status_code" -eq 404 ]]; then
		echo "Snippet with $id not found"
		return
	elif [[ "$status_code" -eq 401 ]]; then
		echo "You are not authorized to access this resource"
		return
	else
		echo "An error occured while deleting the snippet"
		return
	fi
}

quit() {
	echo "Bye for now ..."
	exit 0
}

clear_screen() {
	clear
}

check_dependencies() {
	if command -v jq &>/dev/null command -v curl &>/dev/null && command -v awk &>/dev/null; then
    echo "All required commands are installed."
	else
    echo "One or more required commands are missing. Please install jq, curl, and awk."
    exit 0
	fi
}

main() {
	echo "What do you want to do?" 
	echo "1. Signup"
	echo "2. Signin"
	echo "3. Create Snippet"
	echo "4. Get Snippets"
	echo "5. Get Snippet"
	echo "6. Delete Snippet"
	echo "7. Quit"
	echo "8. Clear"
	read -p "" choice
	case $choice in
		1) signup; main ;;
		2) signin; main ;;
		3) createsnippet; main ;;
		4) getsnippets ; main ;;
		5) getsnippet ; main ;;
		6) deletesnippet ; main ;;
		7) quit ;;
		8) clear_screen ; main ;;
		*) echo "Invalid choice" ; main ;;
	esac
}

check_dependencies

main
