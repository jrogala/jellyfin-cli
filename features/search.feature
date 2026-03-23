Feature: Search
  As a user I can search for media items.

  Background:
    Given a running Jellyfin mock server
    And an authenticated client

  Scenario: Search returns matching items
    Given the server returns search results for "panda":
      | id      | name          | type  | year |
      | mov-001 | Kung Fu Panda | Movie | 2008 |
    When I search for "panda" with limit 10
    Then I should get 1 search results
    And search result "Kung Fu Panda" should have type "Movie"

  Scenario: Search with no results
    Given the server returns no search results for "nonexistent"
    When I search for "nonexistent" with limit 10
    Then I should get 0 search results
