Feature: Item listing
  As a user I can list and inspect media items.

  Background:
    Given a running Jellyfin mock server
    And an authenticated client

  Scenario: List movies
    Given the server has movies:
      | id      | name             | year |
      | mov-001 | Kung Fu Panda    | 2008 |
      | mov-002 | The Matrix       | 1999 |
    When I list movies
    Then I should get 2 items
    And item "Kung Fu Panda" should have year 2008

  Scenario: List items with type filter
    Given the server has items of type "Episode":
      | id      | name       |
      | ep-001  | Pilot      |
    When I list items with type "Episode"
    Then I should get 1 items

  Scenario: Get item info
    Given the server has item "item-42" with name "Interstellar" and year 2014
    When I get info for item "item-42"
    Then the item name should be "Interstellar"
    And the item year should be 2014
