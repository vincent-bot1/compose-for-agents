package catalog

func (c *Catalog) Find(serverName string) (*Server, *map[string]Tool, bool) {
	// Is it an MCP Server?
	server, ok := c.Servers[serverName]
	if ok {
		return &server, nil, true
	}

	// Is it a tool group?
	tools, ok := c.Tools[serverName]
	if ok {
		return nil, &tools, true
	}

	return nil, nil, false
}
