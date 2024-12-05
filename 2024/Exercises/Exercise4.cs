using System.Text.RegularExpressions;
using AdventOfCode.Utilities;

namespace AdventOfCode.Exercises;

public class Exercise4 : IExercise
{
    private int[][] _directions = {
        new int[]{  0,  1 },
        new int[]{  0, -1 },
        new int[]{  1,  0 },
        new int[]{ -1,  0 },
        new int[]{  1,  1 },
        new int[]{  1, -1 },
        new int[]{ -1,  1 },
        new int[]{ -1, -1 },
    };

    private int[][] _diagonalDirections = {
        new int[]{  1,  1 },
        new int[]{  1, -1 },
        new int[]{ -1,  1 },
        new int[]{ -1, -1 },
    };

    private char[] _wordTarget = { 'X', 'M', 'A', 'S' };

    public int GetDay() => 4;

    public Result ExecutePart1(string inputFile)
    {
        var input = File.ReadAllLines(inputFile);
        return new Result(SearchForXmas(input), false);
    }

    private int[] FindChar(string line, char x)
    {
        return line.Select((c, i) => (c, i)).Where(t => t.c == x).Select(t => t.i).ToArray();
    }

    private int[] FindX(string line)
    {
        return FindChar(line, 'X');
    }

    private int[] FindA(string line)
    {
        return FindChar(line, 'A');
    }

    private int CountXmas(string[] data, int row, int col, int[] direction, int wIdx)
    {
        if (wIdx > _wordTarget.Length) return 0;
        if (wIdx == _wordTarget.Length) return 1;

        var nextChar = _wordTarget[wIdx];
        var nextRow = row + direction[0];
        var nextCol = col + direction[1];

        if (nextRow >= data.Length || nextCol >= data[0].Length || nextRow < 0 || nextCol < 0) return 0;

        var newChar = data[nextRow][nextCol];

        if (newChar == nextChar)
        {
            return CountXmas(data, nextRow, nextCol, direction, wIdx + 1);
        }

        return 0;
    }

    private int SearchForXmas(string[] data)
    {
        var counts = 0;
        for (var rowIdx = 0; rowIdx < data.Length; rowIdx++)
        {
            var colIndices = FindX(data[rowIdx]);
            foreach (var colIdx in colIndices)
            {
                foreach (var direction in _directions)
                {
                    counts += CountXmas(data, rowIdx, colIdx, direction, 1);
                }
            }
        }
        return counts;
    }

    private bool BoundsOk(string[] data, int row, int col)
    {
        return (row < data.Length && col < data[0].Length && row >= 0 && col >= 0);
    }

    private int SearchForMAS(string[] data)
    {
        var counts = 0;
        for (var rowIdx = 0; rowIdx < data.Length; rowIdx++)
        {
            var colIndices = FindA(data[rowIdx]);

            foreach (var colIdx in colIndices)
            {
                var mCount = 0;
                var sCount = 0;
                Dictionary<char, List<List<int>>> chars = new() { {'M', []}, {'S', []} };

                foreach (var diag in _diagonalDirections)
                {
                    var nextRow = rowIdx + diag[0];
                    var nextCol = colIdx + diag[1];
                    if (!BoundsOk(data, nextRow, nextCol)) continue;

                    if (data[nextRow][nextCol] == 'M') chars['M'].Add( [nextRow, nextCol] );
                    if (data[nextRow][nextCol] == 'S') chars['S'].Add( [nextRow, nextCol] );
                }

                var mOk = false;
                var sOk = false;
                if (chars['M'].Count() == 2)
                {
                    var f = chars['M'].First();
                    var l = chars['M'].Last();
                    if (f[0] == l[0] || f[1] == l[1])
                    {
                        mOk = true;
                    }
                }
                if (chars['S'].Count() == 2)
                {
                    var f = chars['S'].First();
                    var l = chars['S'].Last();
                    if (f[0] == l[0] || f[1] == l[1])
                    {
                        sOk = true;
                    }
                }

                if (mOk && sOk) counts++;
            }
        }
        return counts;
    }


    public Result ExecutePart2(string inputFile)
    {
        var input = File.ReadAllLines(inputFile);

        return new Result(SearchForMAS(input), false);
    }

}
