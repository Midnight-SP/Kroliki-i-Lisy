# Kroliki-i-Lisy

Symulacja ekosystemu z królikami, lisami i trawą w środowisku 2D z graficzną wizualizacją.

## Opis działania

Program symuluje świat podzielony na kwadratową siatkę, na której:
- **Trawa** rośnie losowo na pustych polach.
- **Króliki** poruszają się, jedzą trawę, rozmnażają się i uciekają przed lisami.
- **Lisy** polują na króliki, rozmnażają się i umierają z głodu.

Każde zwierzę ma energię, która spada z wiekiem i rośnie po zjedzeniu pokarmu. Zwierzęta rozmnażają się, gdy mają wystarczająco dużo energii. Gdy energia spadnie do zera, zwierzę umiera.

## Opis algorytmów

- **Króliki** poruszają się losowo, ale jeśli w pobliżu jest lis, próbują uciekać na najdalsze pole. Jeśli są głodne, szukają trawy. Jeśli mają dużo energii, mogą się rozmnażać.
- **Lisy** szukają królików w sąsiedztwie, a jeśli są najedzone, mogą się rozmnażać. W przeciwnym razie poruszają się losowo.
- **Trawa** rośnie losowo na pustych polach z prawdopodobieństwem określonym przez parametr `GrowthRate`.

## Interfejs użytkownika

1. **Menu startowe**  
   Po uruchomieniu programu pojawia się okno, w którym można ustawić:
   - szerokość i wysokość planszy,
   - liczbę początkowych królików i lisów,
   - tempo wzrostu trawy.

   Nawigacja odbywa się za pomocą klawiatury (strzałki, Enter).

2. **Symulacja**  
   Po zatwierdzeniu parametrów otwiera się okno z wizualizacją świata:
   - Tło, trawa, króliki i lisy są reprezentowane przez tekstury (obrazki PNG).
   - W lewym górnym rogu wyświetlana jest aktualna liczba królików i lisów.
   - Symulacja trwa do momentu zamknięcia okna lub wyginięcia wszystkich zwierząt.

3. **Wykres populacji**  
   Po zakończeniu symulacji automatycznie generowany jest wykres liczby królików i lisów w czasie (`populacje.png`), który otwiera się w domyślnej przeglądarce obrazów.

## Wymagane narzędzia i biblioteki

- **Go** (zalecana wersja 1.18 lub nowsza)
- **[raylib-go](https://github.com/gen2brain/raylib-go)** – do grafiki 2D (instalacja: `go get github.com/gen2brain/raylib-go/raylib`)
- **[gonum/plot](https://github.com/gonum/plot)** – do generowania wykresów (instalacja: `go get gonum.org/v1/plot/...`)
- **Obrazki PNG**: `empty.png`, `grass.png`, `rabbit.png`, `fox.png` w katalogu projektu

## Uruchomienie

1. Upewnij się, że masz zainstalowane Go oraz wymagane biblioteki.
2. Upewnij się że w pliku z programem znajdują się pliki `go.mod`, `go.sum`, oraz `raylib.dll`.
3. Umieść pliki tekstur (`empty.png`, `grass.png`, `rabbit.png`, `fox.png`) w katalogu z programem.
4. Uruchom program:
   ```
   go run main.go
   ```
5. Po zakończeniu symulacji wykres populacji zostanie zapisany jako `populacje.png` i otwarty automatycznie.

## Platformy

Program działa na Windows, Linux i macOS (wymaga Raylib oraz Go).  
Otwieranie wykresu po symulacji jest automatycznie dostosowane do systemu operacyjnego.

## Opis techniczny działania programu

### Struktura kodu

Program został napisany w języku Go i korzysta z biblioteki **raylib-go** do obsługi grafiki oraz **gonum/plot** do generowania wykresów populacji. Cała logika symulacji oraz interfejs użytkownika mieszczą się w jednym pliku źródłowym.

#### Główne elementy programu:

- **Struktura `Cell`** – reprezentuje pojedyncze pole planszy, przechowuje informacje o typie podłoża (trawa/pusto), obecności zwierzęcia, energii, cooldownie rozmnażania i wieku.
- **Struktura `World`** – przechowuje dwuwymiarową tablicę pól (`Grid`), rozmiar planszy oraz parametry symulacji (maksymalna ilość trawy, tempo wzrostu).
- **Menu startowe** – realizowane w Raylib, pozwala ustawić parametry symulacji (rozmiar planszy, liczba zwierząt, tempo wzrostu trawy) za pomocą klawiatury.
- **Pętla symulacji** – działa w osobnej gorutynie, co pozwala na płynne renderowanie i aktualizację stanu świata niezależnie od rysowania okna.
- **Pętla renderująca** – w głównym wątku, odświeża okno Raylib, rysuje planszę i wyświetla liczby zwierząt.
- **Zbieranie danych do wykresu** – po każdej turze do globalnej tablicy zapisywane są liczebności królików i lisów.
- **Generowanie wykresu** – po zakończeniu symulacji (zamknięciu okna lub wymarciu zwierząt) tworzony jest wykres populacji i otwierany w domyślnej przeglądarce obrazów.

### Kluczowe decyzje projektowe

- **Dwuwarstwowa reprezentacja planszy** – każde pole przechowuje osobno informację o podłożu (trawa/pusto) i o zwierzęciu (królik/lis/pusto). Dzięki temu można łatwo obsłużyć sytuacje, gdy na jednym polu jest trawa i zwierzę.
- **Gorutyna do symulacji** – logika symulacji (ruch, jedzenie, rozmnażanie, śmierć) działa w osobnym wątku, a główny wątek zajmuje się tylko rysowaniem. Komunikacja odbywa się przez kanał Go (`chan`), co pozwala na płynne odświeżanie okna i reagowanie na zamknięcie przez użytkownika.
- **Kanały do synchronizacji** – kanały `updateChan` i `quitChan` pozwalają bezpiecznie przekazywać stany świata między gorutyną symulacji a pętlą renderującą oraz obsłużyć zamknięcie okna bez deadlocków.
- **Brak wykresu na żywo** – wykres populacji generowany jest dopiero po zakończeniu symulacji, co upraszcza kod i nie wymaga dynamicznego rysowania wykresu w oknie Raylib.
- **Prosta obsługa menu** – menu startowe jest minimalistyczne, obsługiwane tylko klawiaturą, co pozwala uniknąć zależności od dodatkowych bibliotek GUI.

### Główne funkcje i ich rola

- `ShowMenu()` – wyświetla menu startowe i zwraca wybrane parametry symulacji.
- `NewWorld()` i `Initialize()` – tworzą i losowo rozmieszczają trawę, króliki i lisy na planszy.
- `GrowGrass()`, `MoveRabbits()`, `MoveFoxes()`, `UpdateEnergy()` – realizują logikę wzrostu trawy, ruchu, jedzenia, rozmnażania i śmierci zwierząt.
- `SimulateWithVisualization()` – uruchamia gorutynę symulacji, pętlę renderującą oraz po zakończeniu generuje wykres i otwiera go w przeglądarce.
- `ShowPlot()` – generuje wykres populacji na podstawie zebranych danych.
- `openImage()` – otwiera plik wykresu w domyślnej przeglądarce, niezależnie od systemu operacyjnego.

### Powody takiej architektury

- **Wydzielenie logiki symulacji do osobnej gorutyny** pozwala na płynne działanie interfejsu graficznego i łatwe zatrzymanie symulacji przez użytkownika.
- **Prosta komunikacja przez kanały** minimalizuje ryzyko błędów związanych z równoległością i pozwala na łatwe zakończenie programu.
- **Brak złożonych zależności** – program jest łatwy do uruchomienia i przeniesienia na inne systemy, wymaga tylko Go i Raylib.
- **Minimalistyczny interfejs** – skupia się na funkcjonalności i czytelności kodu, a nie na rozbudowanych efektach graficznych.

## Licencja

Projekt edukacyjny – możesz używać i modyfikować dowolnie.