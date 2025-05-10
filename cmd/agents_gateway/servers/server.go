package servers

import (
	"fmt"
	"strings"
)

func (s *Server) GetTitle() string {
	if s.About.Title != "" {
		return s.About.Title
	}

	return s.Name
}

func (s *Server) GetMaturity() string {
	if strings.HasPrefix(s.Image, "mcpcommunity/") {
		return "Community Image"
	}

	return "Official Image"
}

func (s *Server) GetHubURL() string {
	return fmt.Sprintf("https://hub.docker.com/repository/docker/%s", s.Image)
}

func (s *Server) GetContext() string {
	base := s.Source.Project + ".git"

	if s.GetBranch() != "main" {
		base += "#" + s.Source.Branch
	} else {
		base += "#"
	}

	if s.Source.Directory != "" && s.Source.Directory != "." {
		base += ":" + s.Source.Directory
	}

	return strings.TrimSuffix(base, "#")
}

func (s *Server) GetSourceURL() string {
	source := s.Source.Project + "/tree/" + s.GetBranch()
	if s.Source.Directory != "" {
		source += "/" + s.Source.Directory
	}
	return source
}

func (s *Server) GetBranch() string {
	if s.Source.Branch == "" {
		return "main"
	}
	return s.Source.Branch
}

func (s *Server) GetDockerfileUrl() string {
	base := s.Source.Project + "/blob/" + s.GetBranch()
	if s.Source.Directory != "" {
		base += "/" + s.Source.Directory
	}
	return base + "/" + s.GetDockerfile()
}

func (s *Server) GetDockerfile() string {
	if s.Source.Dockerfile == "" {
		return "Dockerfile"
	}
	return s.Source.Dockerfile
}

func (s *Server) GetReadmeURL() string {
	return s.Source.Project + "/blob/" + s.GetBranch() + "/" + s.About.Readme
}
