Feature: Authentication
  As a user I can authenticate with Jellyfin and get an access token.

  Background:
    Given a running Jellyfin mock server

  Scenario: Successful login
    Given the server accepts credentials "admin" / "password123"
    When I authenticate as "admin" with password "password123"
    Then authentication should succeed
    And I should receive a token
    And I should receive a user ID

  Scenario: Failed login with wrong password
    Given the server rejects all credentials
    When I authenticate as "admin" with password "wrong"
    Then authentication should fail with "authentication failed"
