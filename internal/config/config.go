func (c *Config) IsProduction() bool {
    return c.App.Environment == "production"
}

func (c *Config) IsDevelopment() bool {
    return c.App.Environment == "development"
}

func (c *Config) IsTest() bool {
    return c.App.Environment == "test"
}
