using AdventOfCode.Utilities;

namespace AdventOfCode.Exercises;

public class Exercise8 : IExercise
{
    private record Antenna(int Row, int Col, char Type);

    public int GetDay() => 8;

    private List<Antenna> GetAntennas(string[] input)
    {
        List<Antenna> antennas = new();
        for (var i = 0; i < input.Count(); i++)
        {
            for (var j = 0; j < input[0].Count(); j++)
            {
                if (input[i][j] == '.') continue;
                antennas.Add(new(i, j, input[i][j]));
            }
        }
        return antennas;
    }

    private double Distance(double x1, double y1, double x2, double y2)
    {
        return Math.Sqrt( Math.Pow((x2 - x1) , 2) +  Math.Pow((y2 - y1) , 2) );
    }

    private List<(Antenna, Antenna)> MakeCombinations(List<Antenna> antennas)
    {
        List<(Antenna, Antenna)> combinations = new();
        for (var i = 0; i < antennas.Count(); i++)
        {
            for (var j = i; j < antennas.Count(); j++)
            {
                combinations.Add((antennas[i], antennas[j]));
            }
        }
        return combinations;
    }

    private bool WithinBounds(string[] grid, int row, int col)
    {
        return (row < grid.Count()) && (col < grid[0].Count()) && (row >= 0) && (col >= 0);
    }

    public Result ExecutePart1(string inputFile)
    {
        var input = File.ReadAllLines(inputFile);

        var antennas = GetAntennas(input);

        var antennaGroups = antennas.GroupBy(a => a.Type);

        List<(int, int)> antinodes = new();
        foreach (var g in antennaGroups)
        {
            var combinations = MakeCombinations(g.ToList());

            foreach ( (var a1, var a2) in combinations )
            {
                var d1 = Distance(0, 0, a1.Row, a1.Col);
                var d2 = Distance(0, 0, a2.Row, a2.Col);
                Antenna c, f;
                if (d1 > d2)
                {
                    c = a2;
                    f = a1;
                }
                else
                {
                    c = a1;
                    f = a2;
                }

                var diff = new int[] { f.Row - c.Row, f.Col - c.Col };

                var cAntiRow = c.Row - diff[0];
                var cAntiCol = c.Col - diff[1];
                if (WithinBounds(input, cAntiRow, cAntiCol))
                {
                    if (cAntiRow != a1.Row && cAntiRow != a2.Row && cAntiCol != a1.Col && cAntiCol != a2.Col)
                    {
                        input[cAntiRow] = input[cAntiRow].Remove(cAntiCol, 1).Insert(cAntiCol, "#");
                        antinodes.Add((cAntiRow, cAntiCol));
                    }
                }

                var fAntiRow = f.Row + diff[0];
                var fAntiCol = f.Col + diff[1];
                if (WithinBounds(input, fAntiRow, fAntiCol))
                {
                    if (fAntiRow != a1.Row && fAntiRow != a2.Row && fAntiCol != a1.Col && fAntiCol != a2.Col)
                    {
                        input[fAntiRow] = input[fAntiRow].Remove(fAntiCol, 1).Insert(fAntiCol, "#");
                        antinodes.Add((fAntiRow, fAntiCol));
                    }
                }
                // if (WithinBounds(input, fAntiRow, fAntiCol)) input[fAntiRow] = input[fAntiRow].Remove(fAntiCol, 1).Insert(fAntiCol, "#");
                // if (WithinBounds(input, fAntiRow, fAntiCol)) antinodes.Add((fAntiRow, fAntiCol));
            }
        }

        return new Result(antinodes.Distinct().Count(), true);
    }

    public Result ExecutePart2(string inputFile)
    {
        var input = File.ReadAllLines(inputFile);

        var antennas = GetAntennas(input);

        var antennaGroups = antennas.GroupBy(a => a.Type);

        List<(int, int)> antinodes = new();
        foreach (var g in antennaGroups)
        {
            var combinations = MakeCombinations(g.ToList());

            foreach ( (var a1, var a2) in combinations )
            {
                var d1 = Distance(0, 0, a1.Row, a1.Col);
                var d2 = Distance(0, 0, a2.Row, a2.Col);
                if (d1 ==  d2) continue;
                Antenna c, f;
                if (d1 > d2)
                {
                    c = a2;
                    f = a1;
                }
                else
                {
                    c = a1;
                    f = a2;
                }

                var diff = new int[] { f.Row - c.Row, f.Col - c.Col };

                var cAntiRow = c.Row - diff[0];
                var cAntiCol = c.Col - diff[1];
                while (WithinBounds(input, cAntiRow, cAntiCol))
                {
                    // if (cAntiRow != a1.Row && cAntiRow != a2.Row && cAntiCol != a1.Col && cAntiCol != a2.Col)
                    // {
                        input[cAntiRow] = input[cAntiRow].Remove(cAntiCol, 1).Insert(cAntiCol, "#");
                        antinodes.Add((cAntiRow, cAntiCol));
                    // }
                    cAntiRow = cAntiRow - diff[0];
                    cAntiCol = cAntiCol - diff[1];
                }

                var fAntiRow = f.Row + diff[0];
                var fAntiCol = f.Col + diff[1];
                while (WithinBounds(input, fAntiRow, fAntiCol))
                {
                    // if (fAntiRow != a1.Row && fAntiRow != a2.Row && fAntiCol != a1.Col && fAntiCol != a2.Col)
                    // {
                        input[fAntiRow] = input[fAntiRow].Remove(fAntiCol, 1).Insert(fAntiCol, "#");
                        antinodes.Add((fAntiRow, fAntiCol));
                    // }
                    fAntiRow = fAntiRow + diff[0];
                    fAntiCol = fAntiCol + diff[1];
                }
                // if (WithinBounds(input, fAntiRow, fAntiCol)) input[fAntiRow] = input[fAntiRow].Remove(fAntiCol, 1).Insert(fAntiCol, "#");
                // if (WithinBounds(input, fAntiRow, fAntiCol)) antinodes.Add((fAntiRow, fAntiCol));
                antinodes.Add((a1.Row, a1.Col));
                antinodes.Add((a2.Row, a2.Col));
            }
        }

        // input.PrintLines();

        return new Result(antinodes.Distinct().Count(), false);
    }
}
