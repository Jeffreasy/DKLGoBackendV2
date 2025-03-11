package middleware

// Voeg hier een commentaar toe om te bevestigen dat AuthMiddleware de IAuthMiddleware interface implementeert
// Zorg ervoor dat de compiler dit controleert
var _ IAuthMiddleware = (*AuthMiddleware)(nil)

// Rest van de code blijft hetzelfde...
// ... existing code ...
