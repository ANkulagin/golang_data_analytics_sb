package converter

func (c *Converter) ConvertFile(filePath, srcDir, destDir string) error {
	var content []byte
	_, _, _ = c.splitFrontMatter(content)
	return nil
}

func (c *Converter) splitFrontMatter(content []byte) (*FrontMatter, []byte, error) {
	return &FrontMatter{}, nil, nil
}
