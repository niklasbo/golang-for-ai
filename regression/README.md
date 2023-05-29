# Experiment 1: Regression mit unterschiedlich vielen Graden

## Eingabewerte
- x-Werte
- y-Werte
- Grad des höchsten Polynoms

## Ausgabewerte
- Tabelle mit Faktoren für die Polynome angefangen von x^0

## Erklärung
Lösung des Problems mithilfe einer Vandermonde-Matrix.
Danach Aufruf der Gonum-Funktion `SolveTo`.

## Fazit
Mit Golang und der Bibliothek lassen sich Regressionen mit beliebig hohen Exponenten berechnen.
Jedoch müssen die Algorithmen selbst implementiert werden, in diesem Fall durch eine Umwandlung des Problems in Matrizen-Form.
Eine automatische Güteberechnung R^2 wird nicht angeboten.

## Quellen
https://github.com/gonum/gonum/issues/1759#issuecomment-1005668867
https://pkg.go.dev/gonum.org/v1/gonum@v0.11.0/mat#QR
