package com.putaro.app;
import com.putaro.app.Node;
import junit.framework.Assert;

/**
 * Una instancia de esta clase representa una lista simplemente enlazada cuyos elementos
 * son del tipo Object y cuyas funcionalidades se reducen a consulta de tamaño, inserción 
 * de elementos al principio, remoción y consulta del elemento raíz.
 * 
 * Condiciones de borde: El tamaño de la lista se encuentra entre MIN_LIST_SIZE y MAX_LIST_SIZE
 */
public class SimplifiedLinkedList 
{

	public final static int MAX_LIST_SIZE = 1000;
	
	public final static int MIN_LIST_SIZE = 0;
	
	private Node root = null;
	
	private int size = 0;
		
	/**
	 * @param item Objeto que debe contener el nodo a agregar. 
	 * 
	 * @throws AssertionError En caso que se trate de exceder la cantidad máxima de 
	 * elementos permitidos.
	 */
	public void add(Object item)
	{
		Assert.assertTrue(this.size < MAX_LIST_SIZE);
		try
		{
			this.root.addNextNode(item);
		} catch(NullPointerException e)
		{
			this.root = new Node(item);
		} finally 
		{
			this.size++;
		}
	}
	
	public int getSize()
	{
		return this.size;
	}
	
	/**
	 * @throws AssertionError En caso que se violen las condiciones de borde.
	 * @throws NullPointerException En caso que la lista se encuentre vacía.
	 */
	public void removeFirst() throws Exception
	{	
		try{
			this.root = this.getRoot().getNextNode();
		} catch (NullPointerException e)
		{
			this.root.getItem();
			this.root = null;
		} finally {
			this.size--;
		}
		Assert.assertTrue(this.size >= MIN_LIST_SIZE);
	}
	
	protected Node getRoot()
	{
		return this.root;
	}
	
	/**
	 * @throws AssertionError En caso que se violen las condiciones de borde.
	 * @throws Exception En caso que la lista se encuentre vacía.
	 */
	public Object getFirst() throws Exception 
	{
		Assert.assertTrue(this.size >= MIN_LIST_SIZE);
		Assert.assertTrue(this.size <= MAX_LIST_SIZE);
		return this.root.getItem();
	}
	
	public boolean isEmpty()
	{
		return (this.root == null);
	}
}
