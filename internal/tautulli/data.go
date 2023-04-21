package tautulli

type respGetActivity struct {
	Response struct {
		Data struct {
			Sessions []struct {
				MediaType           string `json:"media_type"`
				ProgressPercent     string `json:"progress_percent"`
				QualityProfile      string `json:"quality_profile"`
				User                string `json:"user"`
				LibraryName         string `json:"library_name"`
				Title               string `json:"title"`
				ParentTitle         string `json:"parent_title"`
				GrandparentTitle    string `json:"grandparent_title"`
				FullTitle           string `json:"full_title"`
				VideoFullResolution string `json:"video_full_resolution"`
				Player              string `json:"player"`
				State               string `json:"state"`
				TranscodeDecision   string `json:"transcode_decision"`
			} `json:"sessions"`
			StreamCountDirectPlay   int `json:"stream_count_direct_play"`
			StreamCountDirectStream int `json:"stream_count_direct_stream"`
			StreamCountTranscode    int `json:"stream_count_transcode"`
			TotalBandwidth          int `json:"total_bandwidth"`
			LanBandwidth            int `json:"lan_bandwidth"`
			WanBandwidth            int `json:"wan_bandwidth"`
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
