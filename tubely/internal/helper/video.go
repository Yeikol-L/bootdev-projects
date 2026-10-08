package helper

import (
	"bytes"
	"encoding/json"
	"errors"
	"os/exec"
)

type MetadataCommandResult struct {
	Streams []struct {
		Index              int    `json:"index"`
		CodecName          string `json:"codec_name,omitempty"`
		CodecLongName      string `json:"codec_long_name,omitempty"`
		Profile            string `json:"profile,omitempty"`
		CodecType          string `json:"codec_type"`
		CodecTagString     string `json:"codec_tag_string"`
		CodecTag           string `json:"codec_tag"`
		MimeCodecString    string `json:"mime_codec_string,omitempty"`
		Width              int    `json:"width,omitempty"`
		Height             int    `json:"height,omitempty"`
		CodedWidth         int    `json:"coded_width,omitempty"`
		CodedHeight        int    `json:"coded_height,omitempty"`
		HasBFrames         int    `json:"has_b_frames,omitempty"`
		SampleAspectRatio  string `json:"sample_aspect_ratio,omitempty"`
		DisplayAspectRatio string `json:"display_aspect_ratio,omitempty"`
		PixFmt             string `json:"pix_fmt,omitempty"`
		Level              int    `json:"level,omitempty"`
		ColorRange         string `json:"color_range,omitempty"`
		ColorSpace         string `json:"color_space,omitempty"`
		ColorTransfer      string `json:"color_transfer,omitempty"`
		ColorPrimaries     string `json:"color_primaries,omitempty"`
		ChromaLocation     string `json:"chroma_location,omitempty"`
		FieldOrder         string `json:"field_order,omitempty"`
		IsAvc              string `json:"is_avc,omitempty"`
		NalLengthSize      string `json:"nal_length_size,omitempty"`
		ID                 string `json:"id"`
		RFrameRate         string `json:"r_frame_rate"`
		AvgFrameRate       string `json:"avg_frame_rate"`
		TimeBase           string `json:"time_base"`
		StartPts           int    `json:"start_pts"`
		StartTime          string `json:"start_time"`
		DurationTs         int    `json:"duration_ts"`
		Duration           string `json:"duration"`
		BitRate            string `json:"bit_rate,omitempty"`
		BitsPerRawSample   string `json:"bits_per_raw_sample,omitempty"`
		NbFrames           string `json:"nb_frames"`
		ExtradataSize      int    `json:"extradata_size"`
		Disposition        struct {
			Default         int `json:"default"`
			Dub             int `json:"dub"`
			Original        int `json:"original"`
			Comment         int `json:"comment"`
			Lyrics          int `json:"lyrics"`
			Karaoke         int `json:"karaoke"`
			Forced          int `json:"forced"`
			HearingImpaired int `json:"hearing_impaired"`
			VisualImpaired  int `json:"visual_impaired"`
			CleanEffects    int `json:"clean_effects"`
			AttachedPic     int `json:"attached_pic"`
			TimedThumbnails int `json:"timed_thumbnails"`
			NonDiegetic     int `json:"non_diegetic"`
			Captions        int `json:"captions"`
			Descriptions    int `json:"descriptions"`
			Metadata        int `json:"metadata"`
			Dependent       int `json:"dependent"`
			StillImage      int `json:"still_image"`
			Multilayer      int `json:"multilayer"`
		} `json:"disposition"`
		Tags struct {
			Language    string `json:"language"`
			HandlerName string `json:"handler_name"`
			Encoder     string `json:"encoder"`
			Timecode    string `json:"timecode"`
		} `json:"tags"`
		SampleFmt      string `json:"sample_fmt,omitempty"`
		SampleRate     string `json:"sample_rate,omitempty"`
		Channels       int    `json:"channels,omitempty"`
		ChannelLayout  string `json:"channel_layout,omitempty"`
		BitsPerSample  int    `json:"bits_per_sample,omitempty"`
		InitialPadding int    `json:"initial_padding,omitempty"`
	} `json:"streams"`
}

func GetVideoAspectRatio(filePath string) (string, error) {
	cmd := exec.Command("ffprobe", "-v", "error", "-print_format", "json", "-show_streams", filePath)
	buffer := new(bytes.Buffer)
	cmd.Stdout = buffer
	err := cmd.Run()
	if err != nil {
		return "", err
	}
	var result MetadataCommandResult
	err = json.Unmarshal(buffer.Bytes(), &result)
	if err != nil {
		return "", err
	}
	if len(result.Streams) == 0 {
		return "", errors.New("Expected at least one stream")
	}
	width := result.Streams[0].Width
	height := result.Streams[0].Height
	n, d := aspect(width, height)
	if n == 16 && d == 9 {
		return "16:9", nil
	}
	if n == 9 && d == 16 {
		return "9:16", nil
	}
	return "other", nil
}
func ProcessVideoFortFastStart(filePath string) (string, error) {
	tempFilePath := filePath + ".procesessing"
	cmd := exec.Command("ffmpeg", "-i", filePath, "-c", "copy", "-movflags", "faststart", "-f", "mp4", tempFilePath)
	err := cmd.Run()
	if err != nil {
		return "", err
	}
	return tempFilePath, nil
}
func aspect(n, d int) (int, int) {
	n1 := n / 16
	if d/n1 == 9 {
		return 16, 9
	}
	n1 = n / 9
	if d/n1 == 16 {
		return 9, 16
	}
	return 0, 0
}

// func gcd(a, b int) int {
// 	for b != 0 {
// 		a, b = b, a%b
// 	}
// 	if a < 0 {
// 		return -a
// 	}
// 	return a
// }

// func aspect(n, d int) (int, int) {
// 	g := gcd(n, d)
// 	if g == 0 {
// 		return n, d // 0/0
// 	}
// 	n, d = n/g, d/g
// 	if d < 0 { // signo en el numerador
// 		n, d = -n, -d
// 	}
// 	return n, d
// }
