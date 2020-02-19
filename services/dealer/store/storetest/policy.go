package storetest

func (suite *Suite) CreatePolicy(s string) {
	suite.Require().NoError(
		suite.Store.CreatePolicy(s),
	)
}

func (suite *Suite) RemovePolicy(s string) {
	suite.Require().NoError(
		suite.Store.RemovePolicy(s),
	)
}

func (suite *Suite) HasPolicy(s string) bool {
	found, err := suite.Store.HasPolicy(s)
	suite.Require().NoError(err)
	return found
}

func (suite *Suite) ListPolicies() []string {
	ls, err := suite.Store.ListPolicies()
	suite.Require().NoError(err)
	return ls
}

func (suite *Suite) TestPolicies() {
	suite.Run("NotFound", func() {
		p := "xxx"
		suite.Require().False(
			suite.HasPolicy(p),
		)

		suite.CreatePolicy(p)
		suite.Require().True(
			suite.Store.HasPolicy(p),
		)
	})

	suite.Run("Idempotence", func() {
		p := "foo"
		suite.Store.CreatePolicy(p)
		suite.Store.CreatePolicy(p)
	})

	suite.Run("Removal", func() {
		p := "to-be-removed"
		suite.CreatePolicy(p)
		suite.RemovePolicy(p)
		suite.Require().False(
			suite.HasPolicy(p),
		)
	})

	suite.Run("List", func() {
		ls := suite.ListPolicies()

		ps := []string{"list-1", "list-2", "list-3"}
		for _, p := range ps {
			suite.CreatePolicy(p)
		}

		suite.Assert().Len(
			suite.ListPolicies(),
			len(ls)+len(ps),
		)
	})
}
