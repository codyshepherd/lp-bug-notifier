package main

type Bug struct {
    BugStruct *BugDevel
    LastChecked string
    Changed bool
}

type BugDevel struct {
	Activity_collection_link                  string
	Attachments_collection_link               string
	Bug_tasks_collection_link                 string
	Bug_watches_collection_link               string
	Cves_collection_link                      string
	Date_created                              string
	Date_last_message                         string
	Date_last_updated                         string
	Date_made_private                         string
	Description                               string
	Duplicate_of_link                         string
	Duplicates_collection_link                string
	Heat                                      string
	Http_etag                                 string
	Id                                        string
	Information_type                          string
	Latest_patch_uploaded                     string
	Linked_branches_collection_link           string
	Linked_merge_proposals_collection_link    string
	Message_count                             string
	Messages_collection_link                  string
	Name                                      string
	Number_of_duplicates                      string
	Other_users_affected_count_with_dupes     string
	Owner_link                                string
	Private                                   string
	Resource_type_link                        string
	Security_related                          string
	Self_link                                 string
	Subscriptions_collection_link             string
	Tags                                      string
	Title                                     string
	Users_affected_collection_link            string
	Users_affected_count                      string
	Users_affected_count_with_dupes           string
	Users_affected_with_dupes_collection_link string
	Users_unaffected_collection_link          string
	Users_unaffected_count                    string
	Web_link                                  string
	Who_made_private_link                     string
}

func NewBug() *Bug {
    var b *Bug = new(Bug)
    b.BugStruct = &BugDevel{}
    b.LastChecked = ""
    b.Changed = false
    return b
}

