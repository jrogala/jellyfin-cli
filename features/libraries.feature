Feature: Library listing
  As a user I can list all libraries on my Jellyfin server.

  Background:
    Given a running Jellyfin mock server
    And an authenticated client

  Scenario: List libraries returns available views
    Given the server has libraries:
      | id     | name   | type    |
      | lib-01 | Movies | movies  |
      | lib-02 | Music  | music   |
    When I list libraries
    Then I should get 2 libraries
    And library "Movies" should have type "movies"

  Scenario: Empty library list
    Given the server has no libraries
    When I list libraries
    Then I should get 0 libraries
