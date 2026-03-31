	const batchSize = 50
	for i := 0; i < len(bulkData); i += batchSize {
		end := i + batchSize
		if end > len(bulkData) {
			end = len(bulkData)
		}

		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		err := l.svcCtx.Model.SensitiveFile.BulkMergeByFileIds(ctx, bulkData[i:end])
		cancel()

		if err != nil {
			l.Errorf("batch %d-%d failed: %v", i, end, err)
			continue
		}
	}