using AdventOfCode.Utilities;

namespace AdventOfCode.Exercises;

public class Exercise1 : IExercise
{
    public int GetDay() => 1;

    public Result ExecutePart1(string inputFile)
    {
        var input = File.ReadAllLines(inputFile);

        List<int> leftList = new();
        List<int> rightList = new();
        foreach (var i in input)
        {
            var split = i.Split(" ");
            int.TryParse(split.First(), out var l);
            leftList.Add(l);
            int.TryParse(split.Last(), out var r);
            rightList.Add(r);
        }

        leftList = QSort(leftList);
        rightList = QSort(rightList);
        leftList.PrintLines();
        rightList.PrintLines();

        List<int> diffs = new();
        for (var i = 0; i < leftList.Count(); i++)
        {
            diffs.Add(Math.Abs(leftList[i] - rightList[i]));
        }

        return new Result(diffs.Sum(), true);
    }

    private List<int> QSort(List<int> l)
    {
        if (l.Count() < 2) return l;
        var lo = 0;
        var hi = l.Count() - 1;

        int p = hi / 2;

        (l[p], l[hi]) = (l[hi], l[p]);

        for (var i = 0; i < l.Count(); i++)
        {
            if (l[i] < l[hi])
            {
                (l[lo], l[i]) = (l[i], l[lo]);
                lo++;
            }
        }

        (l[hi], l[lo]) = (l[lo], l[hi]);

        var midVal = l[lo];
        List<int> outList = new();
        outList.AddRange(QSort(l[0..lo]));
        outList.Add(midVal);
        outList.AddRange(QSort(l[(lo+1)..l.Count()]));

        return outList;
    }

    public Result ExecutePart2(string inputFile)
    {
        var input = File.ReadAllLines(inputFile);

        List<int> leftList = new();
        List<int> rightList = new();
        foreach (var i in input)
        {
            var split = i.Split(" ");
            int.TryParse(split.First(), out var l);
            leftList.Add(l);
            int.TryParse(split.Last(), out var r);
            rightList.Add(r);
        }

        List<int> similarityScore = new();
        foreach (var i in leftList)
        {
            similarityScore.Add(i * rightList.Where(n => n ==  i).Count());
        }

        return new Result(similarityScore.Sum(), true);
    }
}
