# Symulacja Królików i Lisów

To klasyczne zadanie jest ostatnim, które sugeruję, abyście napisali. Pracujemy już nad nim od tygodnia lub dwóch, i po oddaniu projektu "Las", będzie to pewnie dość ciekawe do napisania. Wymagania co do tego zadania są następujące:

## Wymagania

1. **Grafika**:  
    Napisz program z wykorzystaniem grafiki. Możesz użyć bibliotek takich jak Raylib, SDL2, Fyne, Ebiten lub innych popularnych rozwiązań.

2. **Elementy symulacji**:  
    Symulacja powinna zakładać współistnienie trzech elementów:
    - **Trawa**: Może pojawiać się na pustych polach i rosnąć (jej ilość zmienia się od 0 do zadanego maksimum).
    - **Króliki**: Mogą przemieszczać się po planszy (losowo lub "widząc" otoczenie kilku sąsiednich pól). Króliki jedzą trawę, która wtedy znika lub zmniejsza się jej ilość. Najedzone króliki mogą się rozmnażać, tworząc nowego dorosłego królika. Po rozmnożeniu przez jakiś czas nie mogą się ponownie rozmnażać.
    - **Lisy**: Nie jedzą trawy, ale polują na króliki. Lisy mogą poruszać się losowo lub "widząc" otoczenie. Najedzony lis, który spotka innego najedzonego lisa, może stworzyć nowego lisa.

3. **Energia zwierząt**:  
    Zwierzęta tracą energię, jeśli nie jedzą. Gdy energia osiągnie 0, zwierzę umiera i znika z planszy.

4. **Parametry symulacji**:  
    Pozostałe parametry oraz sposób rozwiązania ustal samodzielnie.

5. **Obserwacja populacji**:  
    Symulacja powinna umożliwiać obserwację liczby lisów i królików oraz dynamiki zmian w ich populacji.

## Wymagania techniczne

- Program powinien działać pod Linux bez żadnych specjalnych zabiegów. Sprawdź działanie np. w wirtualnej maszynie, aby upewnić się, że nie wystąpią problemy. Rozwiązania typowo "windowsowe" nie będą akceptowane.

## Funkcjonalności dodatkowe

- **Interfejs graficzny**:  
  W idealnym przypadku program może wykorzystywać bibliotekę Fyne, aby umożliwić początkowe ustawienie parametrów, takich jak:
  - Gęstość populacji królików i lisów.
  - Prędkość wzrostu trawy.
  - Inne parametry symulacji.

- **Widok graficzny**:  
  Program powinien rysować planszę przedstawiającą bieżącą sytuację świata. Nie musisz rysować królików i lisów szczegółowo – wystarczą kropki lub piksele. W dolnej części ekranu powinien znajdować się wykres aktualizowany na bieżąco, pokazujący liczby królików, lisów i trawy oraz dynamikę zmian w populacji.

- **Interakcje użytkownika**:  
  Możesz dodać przyciski do uruchamiania i zatrzymywania symulacji oraz opcję "rysowania" myszą pozycji, w których będą króliki i lisy.

## Do oddania

1. **Kod źródłowy**:  
    Program w wersji źródłowej.

2. **Wersja skompilowana**:  
    Program w wersji zoptymalizowanej, skompilowanej pod Linux.

3. **Dokumentacja**:  
    Opracowanie na temat działania algorytmów, które wykorzystałeś, ciekawych miejsc, Twojej oceny i opinii o symulacji, oraz odkryć, które poczyniłeś podczas pracy nad projektem.

## Uwagi końcowe

Aby lepiej zrozumieć dynamikę takiego systemu oraz czym jest stała reprodukcji, obejrzyj film z Veritasium, który był wcześniej podany jako materiał pomocniczy.

### Przykładowy mockup

Poniżej znajduje się uproszczony mockup projektu. Widać w nim tło z królikami i lisami. W praktyce mogą to być nawet pojedyncze piksele (dla większych rozmiarów planszy). Warto umieścić wykres i przyciski, których kliknięcie uruchomi lub zatrzyma symulację.

