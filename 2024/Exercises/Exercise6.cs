using AdventOfCode.Utilities;

namespace AdventOfCode.Exercises;

public class Exercise6 : IExercise
{
    public int GetDay() => 6;

    private int[][] _directions = new int[][]
    {
        new int[] {-1, 0 },
        new int[] { 0, 1 },
        new int[] { 1, 0 },
        new int[] { 0, -1},
    };
    private char[] _directionChars = new char[] { '^', '>', 'V', '<' };

    private int GetNextDirectionIndex(int idx) => (idx + 1) % 4;

    private int[] FindStartPos(string[] input)
    {
        for (var i = 0; i < input.Count(); i++)
        {
            var row = input[i].ToList();

            foreach (var dir in _directionChars)
            {
                if (row.Contains(dir)) return new int[] { i, row.IndexOf(dir) };
            }
        }
        return new int[2];
    }

    private int[] GetNewPosition(int[] pos, int dirIdx)
    {
        var dirDelta = _directions[dirIdx];
        var newPos = new int[2];
        newPos[0] = dirDelta[0] + pos[0];
        newPos[1] = dirDelta[1] + pos[1];
        return newPos;
    }

    private int CountPositions(string[] input, int[] pos, int directionIdx)
    {
        var visited = input.Select(r => Enumerable.Repeat(false, r.Count()).ToArray()).ToArray();
        visited[pos[0]][pos[1]] = true;
        var count = 1;
        while (true)
        {
            var newPos = GetNewPosition(pos, directionIdx);

            if (newPos[0] >= input.Count() || newPos[1] >= input[0].Count() || newPos[0] < 0 || newPos[1] < 0)
            {
                break;
            }
            else if (input[newPos[0]][newPos[1]] == '#')
            {
                directionIdx = GetNextDirectionIndex(directionIdx);
            }
            else
            {
                pos = newPos;
                count++;
                visited[pos[0]][pos[1]] = true;
                // Console.WriteLine($"MOVED TO: {pos[0]} {pos[1]}");
            }
        }

        return visited.SelectMany(r => r).Where(b => b == true).Count();
    }

    public Result ExecutePart1(string inputFile)
    {
        var input = File.ReadAllLines(inputFile);

        var startPos = FindStartPos(input);
        var directionIdx = 0;

        // Console.WriteLine($"{startPos[0]} {startPos[1]}");

        return new Result(CountPositions(input, startPos, directionIdx), false);
    }


    public Result ExecutePart2(string inputFile)
    {
        var input = File.ReadAllLines(inputFile);

        return new Result();
    }
}
