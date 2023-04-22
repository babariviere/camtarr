package tautulli

type respGetActivity struct {
	Response struct {
		Data struct {
			Sessions []struct {
				FullTitle           string `json:"full_title"`
				GrandparentTitle    string `json:"grandparent_title"`
				LibraryName         string `json:"library_name"`
				MediaType           string `json:"media_type"`
				ParentTitle         string `json:"parent_title"`
				Platform            string `json:"platform"`
				Player              string `json:"player"`
				Product             string `json:"product"`
				ProductVersion      string `json:"product_version"`
				ProgressPercent     string `json:"progress_percent"`
				QualityProfile      string `json:"quality_profile"`
				SessionID           string `json:"session_id"`
				State               string `json:"state"`
				Title               string `json:"title"`
				TranscodeDecision   string `json:"transcode_decision"`
				User                string `json:"user"`
				VideoFullResolution string `json:"video_full_resolution"`
			} `json:"sessions"`
			StreamCountDirectPlay   int    `json:"stream_count_direct_play"`
			StreamCountDirectStream int    `json:"stream_count_direct_stream"`
			StreamCountTranscode    int    `json:"stream_count_transcode"`
			StreamCount             string `json:"stream_count"`
			TotalBandwidth          int    `json:"total_bandwidth"`
			LanBandwidth            int    `json:"lan_bandwidth"`
			WanBandwidth            int    `json:"wan_bandwidth"`
		} `json:"data"`
	} `json:"response"`
}

type respGetHistory struct {
	Response struct {
		Data struct {
			Data []struct {
				Date              int    `json:"date"`
				PlayDuration      int    `json:"play_duration"`
				User              string `json:"user"`
				Product           string `json:"product"`
				Player            string `json:"player"`
				MediaType         string `json:"media_type"`
				FullTitle         string `json:"full_title"`
				Title             string `json:"title"`
				ParentTitle       string `json:"parent_title"`
				GrandparentTitle  string `json:"grandparent_title"`
				TranscodeDecision string `json:"transcode_decision"`
			} `json:"data"`
		} `json:"data"`
	} `json:"response"`
}
