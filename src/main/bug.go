package main

type BugDevel struct {
	activity_collection_link    string
	attachments_collection_link string
	bug_tasks_collection_link   string
	bug_watches_collection_link string
	cves_collection_link        string
	date_created                string
	date_last_message           string
	date_last_updated           string
	date_made_private           string
	description                 string

	duplicate_of_link                      string
	duplicates_collection_link             string
	heat                                   string
	http_etag                              string
	id                                     string
	information_type                       string
	latest_patch_uploaded                  string
	linked_branches_collection_link        string
	linked_merge_proposals_collection_link string
	message_count                          string
	messages_collection_link               string
	name                                   string

	number_of_duplicates                      string
	other_users_affected_count_with_dupes     string
	owner_link                                string
	private                                   string
	resource_type_link                        string
	security_related                          string
	self_link                                 string
	subscriptions_collection_link             string
	tags                                      string
	title                                     string
	users_affected_collection_link            string
	users_affected_count                      string
	users_affected_count_with_dupes           string
	users_affected_with_dupes_collection_link string
	users_unaffected_collection_link          string
	users_unaffected_count                    string
	web_link                                  string
	who_made_private_link                     string
}
