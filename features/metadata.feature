Feature: Metadata management
  As a user I can update, identify, and refresh item metadata.

  Background:
    Given a running Jellyfin mock server
    And an authenticated client

  Scenario: Update item metadata
    Given the server has item "item-42" with name "Old Title" and year 2000
    When I update item "item-42" with name "New Title"
    Then the update should succeed

  Scenario: Identify item with provider IDs
    Given the server accepts identify for item "item-42"
    When I identify item "item-42" as "Interstellar" with IMDB "tt0816692"
    Then the identify should succeed

  Scenario: Refresh metadata
    Given the server accepts refresh for item "item-42"
    When I refresh metadata for item "item-42"
    Then the refresh should succeed

  Scenario: Scan all libraries
    Given the server accepts library scan
    When I scan libraries
    Then the scan should succeed
