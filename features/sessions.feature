Feature: Session listing
  As a user I can view active playback sessions.

  Background:
    Given a running Jellyfin mock server
    And an authenticated client

  Scenario: List active sessions
    Given the server has active sessions:
      | id   | device    | client     | user  | now_playing   |
      | s-01 | Chromecast| Jellyfin Web| alice | Kung Fu Panda |
    When I list sessions
    Then I should get 1 sessions
    And session "s-01" should show "Kung Fu Panda" as now playing

  Scenario: No active sessions
    Given the server has no active sessions
    When I list sessions
    Then I should get 0 sessions
