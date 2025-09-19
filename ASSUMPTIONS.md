# Production System Assumptions

This document outlines the assumptions and improvements that would be implemented in a production environment.

## Event Handling

- **Event Streaming**: For production systems, I would use **Kafka** instead of sending events directly via API calls
  - Provides better reliability and scalability
  - Enables event replay and audit capabilities
  - Supports multiple consumers and decoupled architecture

## Authentication & Security

- **Enhanced Authentication**: Scooter clients and user clients would have more robust authentication mechanisms
  - **Current**: Simple API key authentication (used for development simplicity)
  - **Production**: JWT tokens with proper token management
    - Token expiration and refresh mechanisms
    - Role-based access control (RBAC)
    - Multi-factor authentication for sensitive operations