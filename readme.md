# ShopQL - GraphQL Marketplace Service

A GraphQL-based marketplace backend built with gqlgen.

## Features
- Complete GraphQL API implementation
- In-memory data storage
- Complex resolver relationships:
  - Product ↔ Catalog ↔ Seller relationships
  - Hierarchical catalog structure
- Custom directives support

## Technical Details
- **Schema-first** development approach
- Uses `forceResolver` for complex field resolution
- Proper `input` types for mutations
- Sample data loaded from `testdata.json`

## Development Setup
1. Define GraphQL schema (or use provided)
2. Generate gqlgen boilerplate
3. Implement resolvers and business logic
4. Test with included queries

## Project Structure
- `schema.graphql` - Complete type definitions
- `gqlgen.yml` - Generator configuration
- Resolvers with proper separation of concerns

Note: Implemented without database for GraphQL-focused development.